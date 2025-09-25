package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"role-play-ai/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) (*redis.Client, error) {
	// 解析Redis数据库编号
	db, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		db = 0
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       db,
	})

	// 测试连接
	_, err = client.Ping(Ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	RedisClient = client
	return client, nil
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return RedisClient
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// CacheKey 缓存键生成器
type CacheKey struct {
	Prefix string
	ID     interface{}
}

// String 生成缓存键字符串
func (k *CacheKey) String() string {
	return fmt.Sprintf("%s:%v", k.Prefix, k.ID)
}

// 常用缓存键前缀
const (
	CharacterCachePrefix    = "character"
	ConversationCachePrefix = "conversation"
	UserCachePrefix         = "user"
	SessionCachePrefix      = "session"
	RateLimitPrefix         = "rate_limit"
	AICachePrefix           = "ai_response"
)

// 缓存过期时间
const (
	CharacterCacheExpiry    = 24 * time.Hour     // 角色信息缓存24小时
	ConversationCacheExpiry = 1 * time.Hour      // 对话列表缓存1小时
	UserCacheExpiry         = 30 * time.Minute   // 用户信息缓存30分钟
	SessionCacheExpiry      = 7 * 24 * time.Hour // 会话缓存7天
	RateLimitExpiry         = 1 * time.Minute    // 限流缓存1分钟
	AICacheExpiry           = 1 * time.Hour      // AI响应缓存1小时
)

// SetCache 设置缓存
func SetCache(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(Ctx, key, value, expiration).Err()
}

// GetCache 获取缓存
func GetCache(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}

// DeleteCache 删除缓存
func DeleteCache(key string) error {
	return RedisClient.Del(Ctx, key).Err()
}

// DeleteCachePattern 删除匹配模式的缓存
func DeleteCachePattern(pattern string) error {
	keys, err := RedisClient.Keys(Ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return RedisClient.Del(Ctx, keys...).Err()
	}
	return nil
}

// ExistsCache 检查缓存是否存在
func ExistsCache(key string) (bool, error) {
	result, err := RedisClient.Exists(Ctx, key).Result()
	return result > 0, err
}

// IncrementCache 递增缓存值
func IncrementCache(key string, expiration time.Duration) (int64, error) {
	pipe := RedisClient.Pipeline()
	incr := pipe.Incr(Ctx, key)
	pipe.Expire(Ctx, key, expiration)
	_, err := pipe.Exec(Ctx)
	return incr.Val(), err
}
