package handlers

import (
	"context"
	"followservice/models"
	"followservice/proto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type FollowHandler struct {
	collection        *mongo.Collection
	userServiceClient proto.UserServiceClient
	postServiceClient proto.PostServiceClient
}

func NewFollowHandler(collection *mongo.Collection, userServiceHost, postServiceHost string) (*FollowHandler, error) {
	// 创建用户服务客户端
	userConn, err := grpc.Dial(userServiceHost, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// 创建帖子服务客户端
	postConn, err := grpc.Dial(postServiceHost, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &FollowHandler{
		collection:        collection,
		userServiceClient: proto.NewUserServiceClient(userConn),
		postServiceClient: proto.NewPostServiceClient(postConn),
	}, nil
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

func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	// 获取目标用户ID
	targetUserID := c.Query("targetUserId")
	if len(targetUserID) != 36 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数缺失或格式错误"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户信息"})
		return
	}

	// 检查是否自己取消关注自己
	if userID.(string) == targetUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能取消关注自己"})
		return
	}

	// 检查关注关系是否存在
	exists, err := h.checkFollowExists(c.Request.Context(), userID.(string), targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未关注该用户"})
		return
	}

	// 删除关注关系
	result, err := h.collection.DeleteOne(c.Request.Context(), bson.M{
		"follower_id":  userID.(string),
		"following_id": targetUserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消关注失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "取消关注成功",
	})
}

// GetMyFollowsRequest 定义获取关注列表的请求参数
type GetMyFollowsRequest struct {
	Limit  int `form:"limit,default=10"`
	Offset int `form:"offset,default=0"`
}

// FollowResponse 定义关注列表的响应结构
type FollowResponse struct {
	Follows    []FollowDetail `json:"follows"`
	TotalCount int64          `json:"totalCount"`
}

// FollowDetail 定义每个关注对象的详细信息
type FollowDetail struct {
	TargetUser struct {
		ID       string `json:"id"`
		Avatar   string `json:"avatar"`
		Username string `json:"username"`
	} `json:"targetUser"`
	LatestPostContent string    `json:"latestPostContent"`
	Timestamp         time.Time `json:"timestamp"`
}

// GetMyFollows 获取当前用户的关注列表
func (h *FollowHandler) GetMyFollows(c *gin.Context) {
	var req GetMyFollowsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数缺失或格式错误"})
		return
	}

	// 验证参数
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户信息"})
		return
	}

	// 查询关注列表
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"follower_id": userID.(string),
			},
		},
		{
			"$sort": bson.M{
				"created_at": -1,
			},
		},
		{
			"$skip": req.Offset,
		},
		{
			"$limit": req.Limit,
		},
	}

	cursor, err := h.collection.Aggregate(c.Request.Context(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var follows []models.Follow
	if err := cursor.All(c.Request.Context(), &follows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}

	// 获取总数
	totalCount, err := h.collection.CountDocuments(c.Request.Context(), bson.M{
		"follower_id": userID.(string),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}

	// 构建响应数据
	response := FollowResponse{
		Follows:    make([]FollowDetail, 0, len(follows)),
		TotalCount: totalCount,
	}

	// 获取每个关注用户的详细信息
	for _, follow := range follows {
		// 获取用户信息
		userInfo, err := h.userServiceClient.GetUserInfo(c.Request.Context(), &proto.GetUserInfoRequest{
			UserId: follow.FollowingID,
		})
		if err != nil {
			continue // 跳过获取失败的用户
		}

		// 获取用户最新帖子
		posts, err := h.postServiceClient.GetUserPosts(c.Request.Context(), &proto.GetUserPostsRequest{
			UserId: follow.FollowingID,
			Limit:  1,
			Offset: 0,
		})

		var latestPostContent string
		if err == nil && len(posts.Posts) > 0 {
			latestPostContent = posts.Posts[0].Content
		}

		detail := FollowDetail{
			LatestPostContent: latestPostContent,
			Timestamp:         follow.CreatedAt,
		}
		detail.TargetUser.ID = userInfo.Id
		detail.TargetUser.Avatar = userInfo.Avatar
		detail.TargetUser.Username = userInfo.Username

		response.Follows = append(response.Follows, detail)
	}

	c.JSON(http.StatusOK, response)
}

// GetMyFansRequest 定义获取粉丝列表的请求参数
type GetMyFansRequest struct {
	Limit  int `form:"limit,default=10"`
	Offset int `form:"offset,default=0"`
}

// FansResponse 定义粉丝列表的响应结构
type FansResponse struct {
	Fans       []FanDetail `json:"fans"`
	TotalCount int64       `json:"totalCount"`
}

// FanDetail 定义每个粉丝的详细信息
type FanDetail struct {
	TargetUser struct {
		ID       string `json:"id"`
		Avatar   string `json:"avatar"`
		Username string `json:"username"`
	} `json:"targetUser"`
	LatestPostContent string    `json:"latestPostContent"`
	Timestamp         time.Time `json:"timestamp"`
}

// GetMyFans 获取当前用户的粉丝列表
func (h *FollowHandler) GetMyFans(c *gin.Context) {
	var req GetMyFansRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数缺失或格式错误"})
		return
	}

	// 验证参数
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户信息"})
		return
	}

	// 查询粉丝列表
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"following_id": userID.(string),
			},
		},
		{
			"$sort": bson.M{
				"created_at": -1,
			},
		},
		{
			"$skip": req.Offset,
		},
		{
			"$limit": req.Limit,
		},
	}

	cursor, err := h.collection.Aggregate(c.Request.Context(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var follows []models.Follow
	if err := cursor.All(c.Request.Context(), &follows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}

	// 获取总数
	totalCount, err := h.collection.CountDocuments(c.Request.Context(), bson.M{
		"following_id": userID.(string),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}

	// 构建响应数据
	response := FansResponse{
		Fans:       make([]FanDetail, 0, len(follows)),
		TotalCount: totalCount,
	}

	// 获取每个粉丝的详细信息
	for _, follow := range follows {
		// 获取用户信息
		userInfo, err := h.userServiceClient.GetUserInfo(c.Request.Context(), &proto.GetUserInfoRequest{
			UserId: follow.FollowerID,
		})
		if err != nil {
			continue // 跳过获取失败的用户
		}

		// 获取用户最新帖子
		posts, err := h.postServiceClient.GetUserPosts(c.Request.Context(), &proto.GetUserPostsRequest{
			UserId: follow.FollowerID,
			Limit:  1,
			Offset: 0,
		})

		var latestPostContent string
		if err == nil && len(posts.Posts) > 0 {
			latestPostContent = posts.Posts[0].Content
		}

		detail := FanDetail{
			LatestPostContent: latestPostContent,
			Timestamp:         follow.CreatedAt,
		}
		detail.TargetUser.ID = userInfo.Id
		detail.TargetUser.Avatar = userInfo.Avatar
		detail.TargetUser.Username = userInfo.Username

		response.Fans = append(response.Fans, detail)
	}

	c.JSON(http.StatusOK, response)
}
