package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	SectionID uint   `gorm:"column:section_id;not null;index" json:"sectionId"`
	AuthorID  uint   `gorm:"column:author_id;not null;index" json:"authorId"`
	Title     string `gorm:"column:title;size:255;not null" json:"title"`
	Content   string `gorm:"column:content;type:text" json:"content"`
	Pinned    bool   `gorm:"column:pinned;not null;default:false" json:"pinned"`
	Featured  bool   `gorm:"column:featured;not null;default:false" json:"featured"`
	ViewCount int64  `gorm:"column:view_count;not null;default:0" json:"viewCount"`
}
