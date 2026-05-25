package models

// Like 点赞表
type Like struct {
	ID     int64 `gorm:"column:id;primary_key;auto_increment" json:"id"`
	UserID int64 `gorm:"column:user_id;not null;uniqueIndex:uk_user_post" json:"user_id"`
	PostID int64 `gorm:"column:post_id;not null;uniqueIndex:uk_user_post" json:"post_id"`
}
