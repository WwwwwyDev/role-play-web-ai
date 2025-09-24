package services

import (
	"database/sql"
	"fmt"
	"time"

	"role-play-ai/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(req *models.UserRegister) (*models.User, error) {
	// 检查用户名是否已存在
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("username already exists")
	}

	// 检查邮箱是否已存在
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 插入用户
	result, err := s.db.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		req.Username, req.Email, string(hashedPassword),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// 获取创建的用户
	user, err := s.GetUserByID(int(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get created user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(
		"SELECT id, username, email, created_at, updated_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(
		"SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) VerifyPassword(email, password string) (*models.User, error) {
	var userID int
	var username, emailAddr, passwordHash string
	var createdAt, updatedAt interface{}

	err := s.db.QueryRow(
		"SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = ?",
		email,
	).Scan(&userID, &username, &emailAddr, &passwordHash, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	user := &models.User{
		ID:        userID,
		Username:  username,
		Email:     emailAddr,
		CreatedAt: createdAt.(time.Time),
		UpdatedAt: updatedAt.(time.Time),
	}

	return user, nil
}
