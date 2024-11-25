package middleware

import (
	"context"
	"followservice/proto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type AuthMiddleware struct {
	userClient proto.UserServiceClient
}

func NewAuthMiddleware(userServiceHost string) (*AuthMiddleware, error) {
	conn, err := grpc.Dial(userServiceHost, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{
		userClient: proto.NewUserServiceClient(conn),
	}, nil
}

func (m *AuthMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			return
		}

		// 提取Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的认证格式"})
			return
		}

		// 验证token
		resp, err := m.userClient.ValidateToken(context.Background(), &proto.ValidateTokenRequest{
			Token: tokenParts[1],
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "验证token时发生错误"})
			return
		}

		if !resp.IsValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": resp.Error})
			return
		}

		// 将用户ID存储在上下文中
		c.Set("userId", resp.UserId)
		c.Next()
	}
}
