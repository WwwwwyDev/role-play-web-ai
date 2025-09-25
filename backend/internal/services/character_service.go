package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"role-play-ai/internal/database"
	"role-play-ai/internal/models"
)

type CharacterService struct {
	db *sql.DB
}

func NewCharacterService(db *sql.DB) *CharacterService {
	return &CharacterService{db: db}
}

func (s *CharacterService) GetCharacters() ([]*models.Character, error) {
	// 尝试从Redis缓存获取
	cacheKey := &database.CacheKey{Prefix: database.CharacterCachePrefix, ID: "all"}
	cached, err := database.GetCache(cacheKey.String())
	if err == nil && cached != "" {
		var characters []*models.Character
		if json.Unmarshal([]byte(cached), &characters) == nil {
			return characters, nil
		}
	}

	// 从数据库查询
	rows, err := s.db.Query(`
		SELECT id, name, description, avatar_url, system_prompt, category, created_at, updated_at 
		FROM characters 
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query characters: %w", err)
	}
	defer rows.Close()

	var characters []*models.Character
	for rows.Next() {
		character := &models.Character{}
		err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.Description,
			&character.AvatarURL,
			&character.SystemPrompt,
			&character.Category,
			&character.CreatedAt,
			&character.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}
		characters = append(characters, character)
	}

	// 缓存到Redis
	if data, err := json.Marshal(characters); err == nil {
		database.SetCache(cacheKey.String(), string(data), database.CharacterCacheExpiry)
	}

	return characters, nil
}

func (s *CharacterService) GetCharacter(id int) (*models.Character, error) {
	// 尝试从Redis缓存获取
	cacheKey := &database.CacheKey{Prefix: database.CharacterCachePrefix, ID: id}
	cached, err := database.GetCache(cacheKey.String())
	if err == nil && cached != "" {
		var character models.Character
		if json.Unmarshal([]byte(cached), &character) == nil {
			return &character, nil
		}
	}

	// 从数据库查询
	character := &models.Character{}
	err = s.db.QueryRow(`
		SELECT id, name, description, avatar_url, system_prompt, category, created_at, updated_at 
		FROM characters 
		WHERE id = ?
	`, id).Scan(
		&character.ID,
		&character.Name,
		&character.Description,
		&character.AvatarURL,
		&character.SystemPrompt,
		&character.Category,
		&character.CreatedAt,
		&character.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("character not found")
		}
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	// 缓存到Redis
	if data, err := json.Marshal(character); err == nil {
		database.SetCache(cacheKey.String(), string(data), database.CharacterCacheExpiry)
	}

	return character, nil
}

func (s *CharacterService) SearchCharacters(query string) ([]*models.Character, error) {
	searchTerm := "%" + strings.ToLower(query) + "%"

	rows, err := s.db.Query(`
		SELECT id, name, description, avatar_url, system_prompt, category, created_at, updated_at 
		FROM characters 
		WHERE LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(category) LIKE ?
		ORDER BY name
	`, searchTerm, searchTerm, searchTerm)

	if err != nil {
		return nil, fmt.Errorf("failed to search characters: %w", err)
	}
	defer rows.Close()

	var characters []*models.Character
	for rows.Next() {
		character := &models.Character{}
		err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.Description,
			&character.AvatarURL,
			&character.SystemPrompt,
			&character.Category,
			&character.CreatedAt,
			&character.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}
		characters = append(characters, character)
	}

	return characters, nil
}
