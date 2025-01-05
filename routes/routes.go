package routes

import (
	"field_archive/server/handlers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DefineRoutes(router *gin.Engine, h *handlers.RecordingHandler) {

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "test",
		})
	})

	router.GET("/recordings/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID must be a valid integer",
			})
			return
		}
		h.GetByID(c, id)
	})
}
