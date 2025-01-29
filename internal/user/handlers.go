package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
