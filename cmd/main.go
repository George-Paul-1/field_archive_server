package main

import (
	"context"
	"field_archive/server/internal/config"
	"field_archive/server/internal/database"
	"field_archive/server/repositories"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading Config %v", err)
	}
	// server.Start(cfg, routes.DefineRoutes)
	conn, err := database.Connect(context.Background(), cfg)
	if err != nil {
		fmt.Println(err)
	}
	repo := repositories.NewRecordingRepo(conn)
	recording, err := repo.GetByID(1, context.Background())
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Fetched recording: %+v\n", recording)
	}
}
