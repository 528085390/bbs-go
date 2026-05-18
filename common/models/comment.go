package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID   int64  `gorm:"column:post_id;not null;index" json:"post_id"`
	AuthorID int64  `gorm:"column:author_id;not null;index" json:"author_id"`
	ParentID int64  `gorm:"column:parent_id;index;default:0" json:"parent_id"`
	Content  string `gorm:"type:text;not null" json:"content"`
	Depth    uint32 `gorm:"not null;default:0;check:depth >= 0 AND depth <= 3" json:"depth"`
}
