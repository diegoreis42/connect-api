package auth

import (
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/diegoreis42/connect-api/internal/db"
	"github.com/diegoreis42/connect-api/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func InitParams() *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:       "Connect API",
		Key:         []byte(os.Getenv("JWT_KEY")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: payloadFunc(),

		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	}
}

func HandlerMiddleware(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*user.User); ok {
			return jwt.MapClaims{
				"id":       v.ID,
				"username": v.UserName,
			}
		}
		return jwt.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		id, ok := claims["id"].(float64)
		if !ok {
			return nil
		}
		username, ok := claims["username"].(string)
		if !ok {
			return nil
		}
		return &user.User{
			Model: gorm.Model{
				ID: uint(id),
			},
			UserName: username,
		}
	}
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals login
		if err := c.ShouldBindJSON(&loginVals); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}

		username := loginVals.Username
		password := loginVals.Password

		var foundUser user.User
		if err := db.DB.Where("username = ?", username).First(&foundUser).Error; err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		// Compare hashed password
		if err := bcrypt.CompareHashAndPassword(foundUser.Password, []byte(password)); err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		// Return the user payload
		return &user.User{
			Model:    foundUser.Model, // Includes ID
			UserName: foundUser.UserName,
		}, nil
	}
}


func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func HandleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}
