package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	docs "github.com/diegoreis42/connect-api/docs"
	"github.com/diegoreis42/connect-api/internal/auth"
	"github.com/diegoreis42/connect-api/internal/user"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initializeRoutes(router *gin.Engine, jwtHandle *jwt.GinJWTMiddleware) {
	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath

	// Unauthenticated routes
	v1 := router.Group(basePath)
	{
		v1.POST("/login", jwtHandle.LoginHandler)
		v1.POST("/register", user.RegisterHandler)
	}

	// Authenticated routes
	v1_auth := router.Group(basePath, jwtHandle.MiddlewareFunc())
	{
		v1_auth.GET("/refresh_token", jwtHandle.RefreshHandler)
		v1_auth.PATCH("/user/:user_id/follow", user.FollowUser)
		v1_auth.PATCH("/user/:user_id/unfollow", user.UnfollowUser)
		v1_auth.POST("/post", user.AddPost)
		v1_auth.PATCH("/post/:post_id", user.UpdatePost)
	}

	router.NoRoute(jwtHandle.MiddlewareFunc(), auth.HandleNoRoute())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
