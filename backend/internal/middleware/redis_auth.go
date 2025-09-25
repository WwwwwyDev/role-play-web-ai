package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"role-play-ai/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RedisAuthMiddleware 基于Redis的认证中间件
func RedisAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]

		// 检查Redis中是否存在该token
		sessionKey := &database.CacheKey{Prefix: database.SessionCachePrefix, ID: tokenString}
		exists, err := database.ExistsCache(sessionKey.String())
		if err != nil || !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 解析JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			// 如果token无效，从Redis中删除
			database.DeleteCache(sessionKey.String())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 获取用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("user_id", int(userID))
		c.Set("token", tokenString)

		// 更新token的最后访问时间
		userSessionKey := &database.CacheKey{Prefix: database.SessionCachePrefix, ID: fmt.Sprintf("user:%d", int(userID))}
		sessionData := map[string]interface{}{
			"user_id":    int(userID),
			"token":      tokenString,
			"last_seen":  time.Now().Unix(),
			"created_at": claims["iat"],
		}
		if data, err := json.Marshal(sessionData); err == nil {
			database.SetCache(userSessionKey.String(), string(data), database.SessionCacheExpiry)
		}

		c.Next()
	}
}

// StoreSession 存储用户会话到Redis
func StoreSession(userID int, token string) error {
	// 存储token到会话
	sessionKey := &database.CacheKey{Prefix: database.SessionCachePrefix, ID: token}
	sessionData := map[string]interface{}{
		"user_id":    userID,
		"token":      token,
		"last_seen":  time.Now().Unix(),
		"created_at": time.Now().Unix(),
	}

	if data, err := json.Marshal(sessionData); err == nil {
		return database.SetCache(sessionKey.String(), string(data), database.SessionCacheExpiry)
	}
	return nil
}

// RevokeSession 撤销用户会话
func RevokeSession(token string) error {
	sessionKey := &database.CacheKey{Prefix: database.SessionCachePrefix, ID: token}
	return database.DeleteCache(sessionKey.String())
}

// RevokeAllUserSessions 撤销用户的所有会话
func RevokeAllUserSessions(userID int) error {
	userSessionKey := &database.CacheKey{Prefix: database.SessionCachePrefix, ID: fmt.Sprintf("user:%d", userID)}
	return database.DeleteCache(userSessionKey.String())
}

// GetUserSessions 获取用户的所有活跃会话
func GetUserSessions(userID int) ([]map[string]interface{}, error) {
	userSessionKey := &database.CacheKey{Prefix: database.SessionCachePrefix, ID: fmt.Sprintf("user:%d", userID)}
	cached, err := database.GetCache(userSessionKey.String())
	if err != nil {
		return nil, err
	}

	var sessions []map[string]interface{}
	if cached != "" {
		json.Unmarshal([]byte(cached), &sessions)
	}

	return sessions, nil
}
