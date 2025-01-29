package main

import (
	"log"

	"github.com/diegoreis42/connect-api/internal/db"
	"github.com/diegoreis42/connect-api/internal/user"
	"github.com/diegoreis42/connect-api/router"
)

func main() {
	db.Init()

	err := db.DB.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	router.Initialize()
}
