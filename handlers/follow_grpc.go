package handlers

import (
	"context"
	"followservice/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FollowGrpcServer struct {
	proto.UnimplementedFollowServiceServer
	collection *mongo.Collection
}

func NewFollowGrpcServer(collection *mongo.Collection) *FollowGrpcServer {
	return &FollowGrpcServer{
		collection: collection,
	}
}

func (s *FollowGrpcServer) GetFollowCount(ctx context.Context, req *proto.GetFollowCountRequest) (*proto.GetFollowCountResponse, error) {
	// 获取关注数量
	followingCount, err := s.collection.CountDocuments(ctx, bson.M{
		"follower_id": req.UserId,
	})
	if err != nil {
		return nil, err
	}

	// 获取粉丝数量
	followersCount, err := s.collection.CountDocuments(ctx, bson.M{
		"following_id": req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &proto.GetFollowCountResponse{
		FollowersCount: followersCount,
		FollowingCount: followingCount,
	}, nil
}

func (s *FollowGrpcServer) GetFollowingUserIds(ctx context.Context, req *proto.GetFollowingUserIdsRequest) (*proto.GetFollowingUserIdsResponse, error) {
	// 查询指定用户关注的所有用户ID
	cursor, err := s.collection.Find(ctx, bson.M{
		"follower_id": req.UserId,
	}, &options.FindOptions{
		Projection: bson.M{
			"following_id": 1,
			"_id":          0,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var follows []struct {
		FollowingID string `bson:"following_id"`
	}
	if err := cursor.All(ctx, &follows); err != nil {
		return nil, err
	}

	// 构建响应
	followingIds := make([]string, 0, len(follows))
	for _, follow := range follows {
		followingIds = append(followingIds, follow.FollowingID)
	}

	return &proto.GetFollowingUserIdsResponse{
		FollowingUserIds: followingIds,
	}, nil
}
