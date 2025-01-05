package server

import (
	"field_archive/server/internal/config"
	"field_archive/server/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func Start(cfg *config.Config) {
	log.Print("Starting Server...")
	router := gin.Default()
	routes.DefineRoutes(router)
	router.Run(cfg.Port)
}
