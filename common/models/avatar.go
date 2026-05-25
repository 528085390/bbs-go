package models

type avatar struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint   `gorm:"not null;unique" json:"user_id"`
	URL    string `gorm:"type:varchar(255);not null" json:"url"`
}
