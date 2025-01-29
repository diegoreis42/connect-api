package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName  string
	FirstName string
	Password  string
	Friends   []*User `gorm:"many2many:user_friends"`
}
