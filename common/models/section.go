package models

import "gorm.io/gorm"

type Section struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	OrderIndex  int    `json:"orderIndex"`
	Visibility  bool   `json:"visibility"`
}
