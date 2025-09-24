package main

import (
	"log"
	"os"

	"role-play-ai/internal/config"
	"role-play-ai/internal/database"
	"role-play-ai/internal/handlers"
	"role-play-ai/internal/middleware"
	"role-play-ai/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 初始化服务
	userService := services.NewUserService(db)
	characterService := services.NewCharacterService(db)
	conversationService := services.NewConversationService(db)
	messageService := services.NewMessageService(db)
	aiService := services.NewAIService(cfg)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(userService)
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

	// API路由组
	api := r.Group("/api/v1")

	// 认证路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.GetProfile)
	}

	// 角色路由
	characters := api.Group("/characters")
	{
		characters.GET("/", characterHandler.GetCharacters)
		characters.GET("/search", characterHandler.SearchCharacters)
		characters.GET("/:id", characterHandler.GetCharacter)
	}

	// 对话路由（需要认证）
	conversations := api.Group("/conversations")
	conversations.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		conversations.GET("/", conversationHandler.GetConversations)
		conversations.POST("/", conversationHandler.CreateConversation)
		conversations.GET("/:id", conversationHandler.GetConversation)
		conversations.POST("/:id/messages", conversationHandler.SendMessage)
		conversations.POST("/:id/messages/stream", conversationHandler.SendMessageStream)
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
