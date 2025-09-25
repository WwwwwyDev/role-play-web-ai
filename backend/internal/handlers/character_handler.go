package handlers

import (
	"net/http"
	"strconv"

	"role-play-ai/internal/services"

	"github.com/gin-gonic/gin"
)

type CharacterHandler struct {
	characterService *services.CharacterService
}

func NewCharacterHandler(characterService *services.CharacterService) *CharacterHandler {
	return &CharacterHandler{
		characterService: characterService,
	}
}

// GetCharacters 获取所有角色
// @Summary 获取角色列表
// @Description 获取所有可用的AI角色列表
// @Tags 角色
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "角色列表"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /characters [get]
func (h *CharacterHandler) GetCharacters(c *gin.Context) {
	characters, err := h.characterService.GetCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"characters": characters})
}

// GetCharacter 获取单个角色
// @Summary 获取角色详情
// @Description 根据ID获取特定角色的详细信息
// @Tags 角色
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{} "角色详情"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 404 {object} map[string]string "角色不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /characters/{id} [get]
func (h *CharacterHandler) GetCharacter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	character, err := h.characterService.GetCharacter(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"character": character})
}

// SearchCharacters 搜索角色
// @Summary 搜索角色
// @Description 根据关键词搜索角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param query query string true "搜索关键词"
// @Success 200 {object} map[string]interface{} "搜索结果"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /characters/search [get]
func (h *CharacterHandler) SearchCharacters(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	characters, err := h.characterService.SearchCharacters(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"characters": characters})
}
