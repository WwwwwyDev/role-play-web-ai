// @title AI角色扮演聊天API
// @version 1.0
// @description 基于AI的角色扮演聊天系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT认证，格式：Bearer {token}

package main

import (
	"log"
	"os"

	"role-play-ai/internal/config"
	"role-play-ai/internal/database"
	"role-play-ai/internal/handlers"
	"role-play-ai/internal/middleware"
	"role-play-ai/internal/services"

	_ "role-play-ai/docs" // 导入生成的docs包

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 初始化Redis
	_, err = database.InitRedis(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer database.CloseRedis()

	// 初始化服务
	userService := services.NewUserService(db)
	characterService := services.NewCharacterService(db)
	conversationService := services.NewConversationService(db)
	messageService := services.NewMessageService(db)
	aiService := services.NewAIService(cfg)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(userService, cfg.JWTSecret)
	characterHandler := handlers.NewCharacterHandler(characterService)
	conversationHandler := handlers.NewConversationHandler(conversationService, messageService, aiService)

	// 设置Gin模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组
	api := r.Group("/api/v1")

	// 认证路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", middleware.RedisAuthMiddleware(cfg.JWTSecret), authHandler.Logout)
		auth.GET("/me", middleware.RedisAuthMiddleware(cfg.JWTSecret), authHandler.GetProfile)
	}

	// 角色路由
	characters := api.Group("/characters")
	characters.Use(middleware.APIRateLimit())
	{
		characters.GET("/", characterHandler.GetCharacters)
		characters.GET("/search", characterHandler.SearchCharacters)
		characters.GET("/:id", characterHandler.GetCharacter)
	}

	// 对话路由（需要认证）
	conversations := api.Group("/conversations")
	conversations.Use(middleware.RedisAuthMiddleware(cfg.JWTSecret))
	conversations.Use(middleware.APIRateLimit())
	{
		conversations.GET("/", conversationHandler.GetConversations)
		conversations.POST("/", conversationHandler.CreateConversation)
		conversations.GET("/:id", conversationHandler.GetConversation)
		conversations.POST("/:id/messages", middleware.AIChatRateLimit(), conversationHandler.SendMessage)
		conversations.POST("/:id/messages/stream", middleware.AIChatRateLimit(), conversationHandler.SendMessageStream)
		conversations.DELETE("/:id", conversationHandler.DeleteConversation)
		conversations.DELETE("/batch", conversationHandler.BatchDeleteConversations)
	}

	// 启动服务器
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
