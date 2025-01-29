package main

import (
	"log"
	"net/http"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/diegoreis42/connect-api/internal/auth"
	"github.com/diegoreis42/connect-api/internal/db"
	"github.com/diegoreis42/connect-api/internal/user"
	"github.com/gin-gonic/gin"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {
	db.Init()

	err := db.DB.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	engine := gin.Default()

	authMiddleware, err := jwt.New(auth.InitParams())
	if err != nil {
		log.Fatal("JWT Error: " + err.Error())
	}

	engine.Use(auth.HandlerMiddleware(authMiddleware))
	registerRoute(engine, authMiddleware)

	if err = http.ListenAndServe(":"+port, engine); err != nil {
		log.Fatal(err)
	}
}

func registerRoute(r *gin.Engine, handle *jwt.GinJWTMiddleware) {
	r.POST("/login", handle.LoginHandler)
	r.POST("/register", user.RegisterHandler)
	r.NoRoute(handle.MiddlewareFunc(), auth.HandleNoRoute())

	authRoute := r.Group("/auth", handle.MiddlewareFunc())
	authRoute.GET("/refresh_token", handle.RefreshHandler)
}
