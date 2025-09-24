package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"role-play-ai/internal/models"
)

type ConversationService struct {
	db *sql.DB
}

func NewConversationService(db *sql.DB) *ConversationService {
	return &ConversationService{db: db}
}

func (s *ConversationService) GetConversations(userID int) ([]*models.Conversation, error) {
	rows, err := s.db.Query(`
		SELECT c.id, c.user_id, c.character_id, c.title, c.created_at, c.updated_at,
		       ch.id, ch.name, ch.description, ch.avatar_url, ch.system_prompt, ch.category, ch.created_at, ch.updated_at
		FROM conversations c
		JOIN characters ch ON c.character_id = ch.id
		WHERE c.user_id = ?
		ORDER BY c.updated_at DESC
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to query conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*models.Conversation
	for rows.Next() {
		conversation := &models.Conversation{}
		character := &models.Character{}

		err := rows.Scan(
			&conversation.ID,
			&conversation.UserID,
			&conversation.CharacterID,
			&conversation.Title,
			&conversation.CreatedAt,
			&conversation.UpdatedAt,
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
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}

		conversation.Character = character
		conversations = append(conversations, conversation)
	}

	return conversations, nil
}

func (s *ConversationService) GetConversation(id, userID int) (*models.Conversation, error) {
	conversation := &models.Conversation{}
	character := &models.Character{}

	err := s.db.QueryRow(`
		SELECT c.id, c.user_id, c.character_id, c.title, c.created_at, c.updated_at,
		       ch.id, ch.name, ch.description, ch.avatar_url, ch.system_prompt, ch.category, ch.created_at, ch.updated_at
		FROM conversations c
		JOIN characters ch ON c.character_id = ch.id
		WHERE c.id = ? AND c.user_id = ?
	`, id, userID).Scan(
		&conversation.ID,
		&conversation.UserID,
		&conversation.CharacterID,
		&conversation.Title,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
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
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	conversation.Character = character
	return conversation, nil
}

func (s *ConversationService) CreateConversation(userID int, req *models.CreateConversationRequest) (*models.Conversation, error) {
	// 验证角色是否存在
	var characterName string
	err := s.db.QueryRow("SELECT name FROM characters WHERE id = ?", req.CharacterID).Scan(&characterName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("character not found")
		}
		return nil, fmt.Errorf("failed to verify character: %w", err)
	}

	// 创建对话
	title := fmt.Sprintf("与 %s 的对话", characterName)
	result, err := s.db.Exec(
		"INSERT INTO conversations (user_id, character_id, title) VALUES (?, ?, ?)",
		userID, req.CharacterID, title,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	conversationID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation ID: %w", err)
	}

	// 获取创建的对话
	conversation, err := s.GetConversation(int(conversationID), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created conversation: %w", err)
	}

	return conversation, nil
}

func (s *ConversationService) DeleteConversation(id, userID int) error {
	result, err := s.db.Exec("DELETE FROM conversations WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

func (s *ConversationService) BatchDeleteConversations(ids []int, userID int) (int, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("no conversation IDs provided")
	}

	// 构建批量删除的SQL语句
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)+1)

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	args[len(ids)] = userID

	query := fmt.Sprintf("DELETE FROM conversations WHERE id IN (%s) AND user_id = ?",
		strings.Join(placeholders, ","))

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to batch delete conversations: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

func (s *ConversationService) UpdateConversationTitle(id, userID int, title string) error {
	_, err := s.db.Exec(
		"UPDATE conversations SET title = ?, updated_at = ? WHERE id = ? AND user_id = ?",
		title, time.Now(), id, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update conversation title: %w", err)
	}

	return nil
}
