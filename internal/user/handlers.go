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

// @Summary Add a new post
// @Tags post
// @Accept json
// @Produce json
// @Success 201 {object} Post
// @Router /posts [post]
func AddPost(c *gin.Context) {
	user, err := identityHandler(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := Post{
		Content: input.Content,
		UserID:  user.ID,
	}

	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post": post})
}

// @Summary Update a post
// @Tags post
// @Accept json
// @Produce json
// @Success 200 {object} Post
// @Router /posts/:post_id [patch]
func UpdatePost(c *gin.Context) {
	user, err := identityHandler(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postID := c.Param("post_id")

	var post Post
	if err := db.DB.Where("id = ?", postID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Ensure only the owner can update the post
	if post.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own posts"})
		return
	}

	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Content = input.Content
	if err := db.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully", "post": post})
}

// @Summary Delete a post
// @Tags post
// @Success 200
// @Router /posts/:post_id [delete]
func DeletePost(c *gin.Context) {
	user, err := identityHandler(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postID := c.Param("post_id")

	var post Post
	if err := db.DB.Where("id = ?", postID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Ensure only the owner can delete the post
	if post.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}

	if err := db.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
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
