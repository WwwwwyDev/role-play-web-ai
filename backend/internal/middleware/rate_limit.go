package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"role-play-ai/internal/database"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	MaxRequests int                       // 最大请求数
	Window      time.Duration             // 时间窗口
	KeyFunc     func(*gin.Context) string // 限流键生成函数
}

// DefaultRateLimitConfig 默认限流配置
var DefaultRateLimitConfig = RateLimitConfig{
	MaxRequests: 100, // 每分钟100个请求
	Window:      time.Minute,
	KeyFunc: func(c *gin.Context) string {
		// 默认使用IP地址作为限流键
		return c.ClientIP()
	},
}

// APIRateLimitConfig API限流配置
var APIRateLimitConfig = RateLimitConfig{
	MaxRequests: 60, // 每分钟60个API请求
	Window:      time.Minute,
	KeyFunc: func(c *gin.Context) string {
		// 使用用户ID作为限流键（如果已认证）
		if userID, exists := c.Get("user_id"); exists {
			return fmt.Sprintf("user:%d", userID)
		}
		// 否则使用IP地址
		return c.ClientIP()
	},
}

// AIChatRateLimitConfig AI聊天限流配置
var AIChatRateLimitConfig = RateLimitConfig{
	MaxRequests: 10, // 每分钟10次AI聊天请求
	Window:      time.Minute,
	KeyFunc: func(c *gin.Context) string {
		// 使用用户ID作为限流键
		if userID, exists := c.Get("user_id"); exists {
			return fmt.Sprintf("ai_chat:user:%d", userID)
		}
		return c.ClientIP()
	},
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := config.KeyFunc(c)
		rateLimitKey := &database.CacheKey{Prefix: database.RateLimitPrefix, ID: key}

		// 获取当前请求计数
		currentCount, err := database.IncrementCache(rateLimitKey.String(), config.Window)
		if err != nil {
			// 如果Redis出错，允许请求通过但记录错误
			c.Next()
			return
		}

		// 检查是否超过限制
		if currentCount > int64(config.MaxRequests) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": config.Window.Seconds(),
			})
			c.Abort()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(config.MaxRequests-int(currentCount)))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Window).Unix(), 10))

		c.Next()
	}
}

// DefaultRateLimit 默认限流中间件
func DefaultRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(DefaultRateLimitConfig)
}

// APIRateLimit API限流中间件
func APIRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(APIRateLimitConfig)
}

// AIChatRateLimit AI聊天限流中间件
func AIChatRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(AIChatRateLimitConfig)
}

// GetRateLimitStatus 获取限流状态
func GetRateLimitStatus(key string, config RateLimitConfig) (current int64, limit int, remaining int64, resetTime int64, err error) {
	rateLimitKey := &database.CacheKey{Prefix: database.RateLimitPrefix, ID: key}

	// 获取当前计数
	cached, err := database.GetCache(rateLimitKey.String())
	if err != nil {
		return 0, config.MaxRequests, int64(config.MaxRequests), time.Now().Add(config.Window).Unix(), nil
	}

	current, _ = strconv.ParseInt(cached, 10, 64)
	limit = config.MaxRequests
	remaining = int64(limit) - current
	if remaining < 0 {
		remaining = 0
	}
	resetTime = time.Now().Add(config.Window).Unix()

	return current, limit, remaining, resetTime, nil
}
