package main

import (
	"field_archive/server/internal/config"
	"field_archive/server/internal/server"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading Config %v", err)
	}
	server.Start(cfg)

}
