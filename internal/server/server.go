package server

import (
	"field_archive/server/internal/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(cfg *config.Config, DefineRoutes func(*gin.Engine)) {
	log.Print("Starting Server...")
	router := gin.Default()
	DefineRoutes(router)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting Server!\n%v", err)
	}
}
