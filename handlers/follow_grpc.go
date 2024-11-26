package handlers

import (
	"context"
	"followservice/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
