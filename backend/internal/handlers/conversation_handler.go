package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"role-play-ai/internal/models"
	"role-play-ai/internal/services"

	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	conversationService *services.ConversationService
	messageService      *services.MessageService
	aiService           *services.AIService
}

func NewConversationHandler(
	conversationService *services.ConversationService,
	messageService *services.MessageService,
	aiService *services.AIService,
) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
		messageService:      messageService,
		aiService:           aiService,
	}
}

func (h *ConversationHandler) GetConversations(c *gin.Context) {
	userID := c.GetInt("user_id")

	conversations, err := h.conversationService.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

func (h *ConversationHandler) CreateConversation(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conversation, err := h.conversationService.CreateConversation(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"conversation": conversation})
}

func (h *ConversationHandler) GetConversation(c *gin.Context) {
	userID := c.GetInt("user_id")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	conversation, err := h.conversationService.GetConversation(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 获取消息历史
	messages, err := h.messageService.GetConversationMessages(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversation": conversation,
		"messages":     messages,
	})
}

func (h *ConversationHandler) SendMessage(c *gin.Context) {
	userID := c.GetInt("user_id")

	idStr := c.Param("id")
	conversationID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取对话信息
	conversation, err := h.conversationService.GetConversation(conversationID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 保存用户消息
	userMessage, err := h.messageService.CreateMessage(conversationID, "user", req.Content, req.AudioURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取对话历史
	messages, err := h.messageService.GetMessages(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 生成AI响应
	aiResponse, err := h.aiService.GenerateResponse(conversation.Character, messages)
	if err != nil {
		aiResponse = "系统错误，请稍后再试"
	}

	// 保存AI消息
	aiMessage, err := h.messageService.CreateMessage(conversationID, "assistant", aiResponse, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_message": userMessage,
		"ai_message":   aiMessage,
	})
}

// SendMessageStream 发送消息并返回流式AI响应
func (h *ConversationHandler) SendMessageStream(c *gin.Context) {
	userID := c.GetInt("user_id")
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// 验证对话所有权
	conversation, err := h.conversationService.GetConversation(conversationID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	var request struct {
		Content  string `json:"content" binding:"required"`
		AudioURL string `json:"audio_url,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 保存用户消息
	userMessage, err := h.messageService.CreateMessage(conversationID, "user", request.Content, &request.AudioURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取对话历史
	messages, err := h.messageService.GetMessages(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 首先发送用户消息
	userMessageJSON, _ := json.Marshal(userMessage)
	fmt.Fprintf(c.Writer, "data: %s\n\n", string(userMessageJSON))
	c.Writer.Flush()

	// 使用流式响应生成AI回复
	fullResponse, err := h.aiService.GenerateStreamingResponse(conversation.Character, messages, c.Writer)
	if err != nil {
		// 如果流式响应失败，创建错误消息
		errorMsg := "系统错误，请稍后再试"
		aiMessage, _ := h.messageService.CreateMessage(conversationID, "assistant", errorMsg, nil)
		aiMessageJSON, _ := json.Marshal(aiMessage)
		fmt.Fprintf(c.Writer, "data: %s\n\n", string(aiMessageJSON))
		c.Writer.Flush()
		return
	}

	// 流式响应完成后，创建完整的AI消息并保存到数据库
	if fullResponse != "" {
		aiMessage, err := h.messageService.CreateMessage(conversationID, "assistant", fullResponse, nil)
		if err != nil {
			fmt.Printf("Failed to save AI message to database: %v\n", err)
		} else {
			fmt.Printf("AI message saved to database with ID: %d\n", aiMessage.ID)
		}
	}
}

func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
	userID := c.GetInt("user_id")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	err = h.conversationService.DeleteConversation(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted successfully"})
}

func (h *ConversationHandler) BatchDeleteConversations(c *gin.Context) {
	userID := c.GetInt("user_id")

	var request struct {
		IDs []int `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(request.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No conversation IDs provided"})
		return
	}

	deletedCount, err := h.conversationService.BatchDeleteConversations(request.IDs, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Conversations deleted successfully",
		"deleted_count": deletedCount,
	})
}
