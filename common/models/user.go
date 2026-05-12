package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Email       string   `json:"email"`
	DisplayName string   `json:"display_name"`
	Roles       []string `gorm:"serializer:json;type:json"`
}

func NewUser(username, password, email string) *User {
	return &User{
		Username:    username,
		Password:    password,
		Email:       email,
		Roles:       []string{"user"},
		DisplayName: username,
	}
}
