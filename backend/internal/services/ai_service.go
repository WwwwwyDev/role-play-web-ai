package services

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"role-play-ai/internal/config"
	"role-play-ai/internal/database"
	"role-play-ai/internal/models"
)

type AIService struct {
	baseURL string
	model   string
	client  *http.Client
}

type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		baseURL: cfg.OllamaBaseURL,
		model:   cfg.OllamaModel,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *AIService) GenerateResponse(character *models.Character, messages []*models.Message) (string, error) {
	// 生成缓存键
	cacheKey := s.generateCacheKey(character, messages)

	// 尝试从缓存获取响应
	cached, err := database.GetCache(cacheKey)
	if err == nil && cached != "" {
		return cached, nil
	}

	// 构建消息历史
	ollamaMessages := []Message{
		{
			Role:    "system",
			Content: character.SystemPrompt,
		},
	}

	// 添加对话历史
	for _, msg := range messages {
		ollamaMessages = append(ollamaMessages, Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 构建请求
	request := OllamaRequest{
		Model:    s.model,
		Messages: ollamaMessages,
		Stream:   false,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 发送请求到Ollama
	resp, err := s.client.Post(
		fmt.Sprintf("%s/api/chat", s.baseURL),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// 缓存响应
	database.SetCache(cacheKey, ollamaResp.Message.Content, database.AICacheExpiry)

	return ollamaResp.Message.Content, nil
}

// GenerateStreamingResponse 生成流式响应
func (s *AIService) GenerateStreamingResponse(character *models.Character, messages []*models.Message, writer io.Writer) (string, error) {
	// 构建消息历史
	ollamaMessages := []Message{
		{
			Role:    "system",
			Content: character.SystemPrompt,
		},
	}

	// 添加对话历史
	for _, msg := range messages {
		ollamaMessages = append(ollamaMessages, Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 构建请求
	request := OllamaRequest{
		Model:    s.model,
		Messages: ollamaMessages,
		Stream:   true, // 启用流式响应
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 发送请求到Ollama
	resp, err := s.client.Post(
		fmt.Sprintf("%s/api/chat", s.baseURL),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	// 处理流式响应
	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder
	var aiMessageID int64 = 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 解析每一行的JSON响应
		var streamResp OllamaResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			continue // 跳过无效的JSON行
		}

		// 发送内容到writer
		if streamResp.Message.Content != "" {
			// 累积完整响应
			fullResponse.WriteString(streamResp.Message.Content)

			// 创建AI消息对象（第一次时创建，后续更新）
			if aiMessageID == 0 {
				// 创建临时的AI消息对象
				aiMessage := map[string]interface{}{
					"id":              time.Now().UnixNano(),
					"role":            "assistant",
					"content":         streamResp.Message.Content,
					"created_at":      time.Now().Format(time.RFC3339),
					"conversation_id": 0, // 将在后续更新
				}
				aiMessageID = aiMessage["id"].(int64)

				// 发送完整的AI消息对象
				aiMessageJSON, _ := json.Marshal(aiMessage)
				fmt.Fprintf(writer, "data: %s\n\n", string(aiMessageJSON))
			} else {
				// 更新AI消息内容
				aiMessage := map[string]interface{}{
					"id":              aiMessageID,
					"role":            "assistant",
					"content":         fullResponse.String(),
					"created_at":      time.Now().Format(time.RFC3339),
					"conversation_id": 0, // 将在后续更新
				}

				// 发送更新的AI消息对象
				aiMessageJSON, _ := json.Marshal(aiMessage)
				fmt.Fprintf(writer, "data: %s\n\n", string(aiMessageJSON))
			}

			if f, ok := writer.(http.Flusher); ok {
				f.Flush()
			}
		}

		// 如果完成，发送结束标记
		if streamResp.Done {
			fmt.Fprintf(writer, "data: [DONE]\n\n")
			if f, ok := writer.(http.Flusher); ok {
				f.Flush()
			}
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return fullResponse.String(), nil
}

func (s *AIService) TestConnection() error {
	// 测试Ollama连接
	resp, err := s.client.Get(fmt.Sprintf("%s/api/tags", s.baseURL))
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama is not responding properly")
	}

	return nil
}

// generateCacheKey 生成AI响应的缓存键
func (s *AIService) generateCacheKey(character *models.Character, messages []*models.Message) string {
	// 构建缓存数据
	cacheData := map[string]interface{}{
		"character_id": character.ID,
		"model":        s.model,
		"messages":     messages,
	}

	// 序列化为JSON
	data, err := json.Marshal(cacheData)
	if err != nil {
		// 如果序列化失败，使用简单字符串
		return fmt.Sprintf("%s:%d:%s", database.AICachePrefix, character.ID, time.Now().Format("20060102150405"))
	}

	// 生成MD5哈希
	hash := md5.Sum(data)
	return fmt.Sprintf("%s:%s", database.AICachePrefix, hex.EncodeToString(hash[:]))
}

// ClearAICache 清除AI缓存
func (s *AIService) ClearAICache(characterID int) error {
	pattern := fmt.Sprintf("%s:*", database.AICachePrefix)
	return database.DeleteCachePattern(pattern)
}

// GetCacheStats 获取缓存统计信息
func (s *AIService) GetCacheStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取AI缓存数量
	aiKeys, err := database.RedisClient.Keys(database.Ctx, fmt.Sprintf("%s:*", database.AICachePrefix)).Result()
	if err != nil {
		return nil, err
	}
	stats["ai_cache_count"] = len(aiKeys)

	// 获取角色缓存数量
	characterKeys, err := database.RedisClient.Keys(database.Ctx, fmt.Sprintf("%s:*", database.CharacterCachePrefix)).Result()
	if err != nil {
		return nil, err
	}
	stats["character_cache_count"] = len(characterKeys)

	// 获取对话缓存数量
	conversationKeys, err := database.RedisClient.Keys(database.Ctx, fmt.Sprintf("%s:*", database.ConversationCachePrefix)).Result()
	if err != nil {
		return nil, err
	}
	stats["conversation_cache_count"] = len(conversationKeys)

	return stats, nil
}
