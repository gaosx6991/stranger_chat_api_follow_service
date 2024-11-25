package handlers

import (
	"context"
	"followservice/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FollowHandler struct {
	collection *mongo.Collection
}

func NewFollowHandler(collection *mongo.Collection) *FollowHandler {
	return &FollowHandler{
		collection: collection,
	}
}

type FollowUserRequest struct {
	TargetUserID string `json:"targetUserId" binding:"required,len=36"`
}

func (h *FollowHandler) FollowUser(c *gin.Context) {
	var req FollowUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户信息"})
		return
	}

	// 检查是否自己关注自己
	if userID.(string) == req.TargetUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能关注自己"})
		return
	}

	// 检查是否已经关注
	exists, err := h.checkFollowExists(c.Request.Context(), userID.(string), req.TargetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已经关注该用户"})
		return
	}

	// 创建关注关系
	follow := models.Follow{
		ID:          uuid.New().String(),
		FollowerID:  userID.(string),
		FollowingID: req.TargetUserID,
		CreatedAt:   time.Now(),
	}

	_, err = h.collection.InsertOne(c.Request.Context(), follow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "关注成功",
	})
}

func (h *FollowHandler) checkFollowExists(ctx context.Context, followerID, followingID string) (bool, error) {
	count, err := h.collection.CountDocuments(ctx, bson.M{
		"follower_id":  followerID,
		"following_id": followingID,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
