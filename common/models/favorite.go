package models

import (
	"time"
)

// Favorite 收藏表
type Favorite struct {
	ID        int64     `gorm:"column:id;primary_key;auto_increment" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null;uniqueIndex:uk_user_post" json:"user_id"`
	PostID    int64     `gorm:"column:post_id;not null;uniqueIndex:uk_user_post" json:"post_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
