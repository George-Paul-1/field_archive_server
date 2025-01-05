package main

import (
	"field_archive/server/internal/config"
	"field_archive/server/internal/server"
	"field_archive/server/routes"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading Config %v", err)
	}
	server.Start(cfg, routes.DefineRoutes)
}
