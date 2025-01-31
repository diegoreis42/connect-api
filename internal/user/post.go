package user

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Content string `gorm:"column:content"`
	UserID  uint
}

type PostInput struct {
	Content string `json:"content" binding:"required"`
}
