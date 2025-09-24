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

func (h *CharacterHandler) GetCharacters(c *gin.Context) {
	characters, err := h.characterService.GetCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"characters": characters})
}

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
