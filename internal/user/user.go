package user

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/diegoreis42/connect-api/internal/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName  string  `gorm:"column:username;unique"`
	FirstName string  `gorm:"column:first_name"`
	Password  []byte  `json:"-"`
	Friends   []*User `gorm:"many2many:user_friends"`
}

type UserInput struct {
	UserName  string `json:"username" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func RegisterHandler(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := register(input)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": gin.H{
		"id":         user.ID,
		"username":   user.UserName,
		"first_name": user.FirstName,
	}})
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
		Friends:   []*User{},
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return User{}, fmt.Errorf("error while creating user: %v", result.Error)
	}

	return user, nil
}
