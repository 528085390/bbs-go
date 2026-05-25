package models

import "time"

type File struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ObjectKey string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"object_key"`
	Filename  string    `gorm:"type:varchar(255);not null" json:"filename"`
	UserId    int64     `gorm:"not null" json:"user_id"`
	PostId    int64     `gorm:"default:0" json:"post_id"`
	Type      string    `gorm:"type:varchar(32);default:'other'" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
