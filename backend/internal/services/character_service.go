package services

import (
	"database/sql"
	"fmt"
	"strings"

	"role-play-ai/internal/models"
)

type CharacterService struct {
	db *sql.DB
}

func NewCharacterService(db *sql.DB) *CharacterService {
	return &CharacterService{db: db}
}

func (s *CharacterService) GetCharacters() ([]*models.Character, error) {
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

	return characters, nil
}

func (s *CharacterService) GetCharacter(id int) (*models.Character, error) {
	character := &models.Character{}
	err := s.db.QueryRow(`
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
