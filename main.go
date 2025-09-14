package main

import (
	"log"
	"os"

	"your/module/path/config"
	"your/module/path/repositories"
	"your/module/path/routes"
)

func main() {
	if err := config.InitDB(); err != nil {
		log.Fatalf("DB init error: %v", err)
	}

	repo := repositories.NewTodoRepo()
	if err := repo.AutoMigrate(); err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}

	r := routes.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
