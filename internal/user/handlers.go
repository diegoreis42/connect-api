package user

import (
	"errors"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/diegoreis42/connect-api/internal/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @BasePath /api/v1

// @Summary Register User
// @Tags user
// @Accept json
// @Produce json
// @Sucess 201 {object} UserSchema
// @Router /register [post]
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

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"username":   user.UserName,
		"first_name": user.FirstName,
	})
}

// @Summary Follow new User
// @Tags user
// @Sucess 200
// @Router /user/:user_id/follow [patch]
func FollowUser(c *gin.Context) {
	_, err := identityHandler(c)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	}

	user_id := c.Param("user_id")

	var user, follower User

	if err := db.DB.Where("id = ?", user_id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User to follow not found"})
		return
	}

	db.DB.Model(&follower).Association("Following").Append(&user)
	db.DB.Model(&user).Association("Followers").Append(&follower)

	c.JSON(http.StatusOK, gin.H{"message": "User followed successfully"})
}

// @Summary Unfollow an User
// @Tags user
// @Sucess 200
// @Router /user/:user_id/unfollow [patch]
func UnfollowUser(c *gin.Context) {
	_, err := identityHandler(c)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	}

	user_id := c.Param("user_id")

	var user, follower User

	if err := db.DB.Where("id = ?", user_id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User to unfollow not found"})
		return
	}

	db.DB.Model(&follower).Association("Following").Delete(&user)
	db.DB.Model(&user).Association("Followers").Delete(&follower)

	c.JSON(http.StatusOK, gin.H{"message": "User unfollowed successfully"})
}

func identityHandler(c *gin.Context) (User, error) {

	claims := jwt.ExtractClaims(c)
	id, ok := claims["id"].(float64)
	if !ok {
		return User{}, errors.New("error getting id from JWT")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return User{}, errors.New("error getting username from JWT")
	}
	return User{
		Model: gorm.Model{
			ID: uint(id),
		},
		UserName: username,
	}, nil
}
