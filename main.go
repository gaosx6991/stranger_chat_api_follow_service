package main

import (
	"context"
	"fmt"
	"followservice/config"
	"followservice/handlers"
	"followservice/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 连接MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatalf("无法连接MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	collection := mongoClient.Database(cfg.MongoDB.Database).Collection(cfg.MongoDB.Collection)

	// 创建认证中间件
	authMiddleware, err := middleware.NewAuthMiddleware(cfg.UserService.Host)
	if err != nil {
		log.Fatalf("无法创建认证中间件: %v", err)
	}

	// 创建处理器
	followHandler, err := handlers.NewFollowHandler(
		collection,
		cfg.UserService.Host,
		cfg.PostService.Host,
	)
	if err != nil {
		log.Fatalf("无法创建处理器: %v", err)
	}

	// 设置路由
	r := gin.Default()

	// API路由组
	api := r.Group("/api/v1")
	{
		follow := api.Group("/follow")
		{
			follow.POST("/user", authMiddleware.ValidateToken(), followHandler.FollowUser)
			follow.DELETE("/user", authMiddleware.ValidateToken(), followHandler.UnfollowUser)
			follow.GET("/my-follows", authMiddleware.ValidateToken(), followHandler.GetMyFollows)
		}
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
