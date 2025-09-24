package models

import "time"

// User 用户模型
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRegister 用户注册请求
type UserRegister struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserLogin 用户登录请求
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Character 角色模型
type Character struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	AvatarURL    string    `json:"avatar_url" db:"avatar_url"`
	SystemPrompt string    `json:"system_prompt" db:"system_prompt"`
	Category     string    `json:"category" db:"category"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Conversation 对话会话模型
type Conversation struct {
	ID          int        `json:"id" db:"id"`
	UserID      int        `json:"user_id" db:"user_id"`
	CharacterID int        `json:"character_id" db:"character_id"`
	Title       string     `json:"title" db:"title"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	Character   *Character `json:"character,omitempty"`
}

// Message 消息模型
type Message struct {
	ID             int       `json:"id" db:"id"`
	ConversationID int       `json:"conversation_id" db:"conversation_id"`
	Role           string    `json:"role" db:"role"`
	Content        string    `json:"content" db:"content"`
	AudioURL       *string   `json:"audio_url,omitempty" db:"audio_url"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// CreateConversationRequest 创建对话请求
type CreateConversationRequest struct {
	CharacterID int `json:"character_id" binding:"required"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	Content  string  `json:"content" binding:"required"`
	AudioURL *string `json:"audio_url,omitempty"`
}

// AIResponse AI响应
type AIResponse struct {
	Content string `json:"content"`
}
