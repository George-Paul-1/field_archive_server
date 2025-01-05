package routes

import (
	"field_archive/server/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefineRoutes(router *gin.Engine, h *handlers.RecordingHandler) {

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "test",
		})
	})

	router.GET("/recordings/:id", func(c *gin.Context) {
		h.GetByID(c)
	})
}
