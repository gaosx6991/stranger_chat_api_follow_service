package models

import (
	"time"
)

type Follow struct {
	ID          string    `bson:"_id"`
	FollowerID  string    `bson:"follower_id"`
	FollowingID string    `bson:"following_id"`
	CreatedAt   time.Time `bson:"created_at"`
}
