package main

import (
	"context"
	"field_archive/server/handlers"
	"field_archive/server/internal/config"
	"field_archive/server/internal/database"
	"field_archive/server/internal/server"
	"field_archive/server/repositories"
	"field_archive/server/routes"
	"field_archive/server/services"
	"log"
)

func main() {
	// Loading configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading Config %v", err)
	}

	// Building database connection
	db, err := database.Connect(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Couldn't connect to Database %v", err)
	}

	// Establishing 'recordings' interactors
	repo := repositories.NewRecordingRepo(db)
	service := services.NewRecordingService(repo)
	handler := handlers.NewRecordingHandler(service)

	// Starting server
	server.Start(cfg, routes.DefineRoutes, handler)
}
