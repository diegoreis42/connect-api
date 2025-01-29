package router

import (
	"log"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/diegoreis42/connect-api/internal/auth"
	"github.com/gin-gonic/gin"
)

func Initialize() {

	engine := gin.Default()
	authMiddleware, err := jwt.New(auth.InitParams())
	if err != nil {
		log.Fatal("JWT Error: " + err.Error())
	}

	engine.Use(auth.HandlerMiddleware(authMiddleware))

	initializeRoutes(engine, authMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	engine.Run("0.0.0.0:" + port)
}
