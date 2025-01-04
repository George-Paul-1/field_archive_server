package server

import (
	"field_archive/server/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(cfg *config.Config) {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	router.Run(cfg.Port)
}
