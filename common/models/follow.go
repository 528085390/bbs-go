package models

import "time"

// Follow 关注表
type Follow struct {
	ID          int64     `gorm:"column:id;primary_key;auto_increment" json:"id"`
	FollowerID  int64     `gorm:"column:follower_id;not null;uniqueIndex:uk_follower_following" json:"follower_id"`
	FollowingID int64     `gorm:"column:following_id;not null;uniqueIndex:uk_follower_following" json:"following_id"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}
