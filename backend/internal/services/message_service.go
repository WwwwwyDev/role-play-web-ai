package services

import (
	"database/sql"
	"fmt"

	"role-play-ai/internal/models"
)

type MessageService struct {
	db *sql.DB
}

func NewMessageService(db *sql.DB) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) GetMessages(conversationID int) ([]*models.Message, error) {
	rows, err := s.db.Query(`
		SELECT id, conversation_id, role, content, audio_url, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at ASC
	`, conversationID)

	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.Role,
			&message.Content,
			&message.AudioURL,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (s *MessageService) CreateMessage(conversationID int, role, content string, audioURL *string) (*models.Message, error) {
	result, err := s.db.Exec(
		"INSERT INTO messages (conversation_id, role, content, audio_url) VALUES (?, ?, ?, ?)",
		conversationID, role, content, audioURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get message ID: %w", err)
	}

	// 获取创建的消息
	message, err := s.GetMessage(int(messageID))
	if err != nil {
		return nil, fmt.Errorf("failed to get created message: %w", err)
	}

	return message, nil
}

func (s *MessageService) GetMessage(id int) (*models.Message, error) {
	message := &models.Message{}
	err := s.db.QueryRow(`
		SELECT id, conversation_id, role, content, audio_url, created_at
		FROM messages
		WHERE id = ?
	`, id).Scan(
		&message.ID,
		&message.ConversationID,
		&message.Role,
		&message.Content,
		&message.AudioURL,
		&message.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return message, nil
}

func (s *MessageService) GetConversationMessages(conversationID, userID int) ([]*models.Message, error) {
	// 验证对话是否属于用户
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM conversations WHERE id = ? AND user_id = ?", conversationID, userID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to verify conversation: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("conversation not found")
	}

	return s.GetMessages(conversationID)
}

// UpdateMessage 更新消息内容
func (s *MessageService) UpdateMessage(messageID int, content string) error {
	_, err := s.db.Exec("UPDATE messages SET content = ? WHERE id = ?", content, messageID)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	return nil
}
