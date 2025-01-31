package user

import (
	"fmt"
	"os"
	"strconv"

	"github.com/diegoreis42/connect-api/internal/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName  string  `gorm:"column:username;unique"`
	FirstName string  `gorm:"column:first_name"`
	Password  []byte  `json:"-"`
	Followers []*User `gorm:"many2many:user_followers"`
	Following []*User `gorm:"many2many:following_users"`
	Posts     []Post  `gorm:"foreignKey:UserID"`
}

type UserInput struct {
	UserName  string `json:"username" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func register(input UserInput) (User, error) {

	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		cost = 12
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), cost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %v", err)
	}

	user := User{
		UserName:  input.UserName,
		FirstName: input.FirstName,
		Password:  hashedPassword,
		Following: []*User{},
		Followers: []*User{},
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return User{}, fmt.Errorf("error while creating user: %v", result.Error)
	}

	return user, nil
}
