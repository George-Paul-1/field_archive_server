package routes

import (
	"field_archive/server/handlers"
	"field_archive/server/internal/config"
	"field_archive/server/internal/database"
	"field_archive/server/repositories"
	"net/http"
	"os"

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

	router.GET("/recordings/list/:limit", func(c *gin.Context) {
		h.ListItems(c)
	})

	router.GET("/audio/*filepath", func(c *gin.Context) {

		// TODO shift this code to handlers package
		path := c.Param("filepath")

		if _, err := os.Stat(path); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.File(path)
	})

	// TEST ROUTE FOR JS MAP

	router.GET("/locations", func(c *gin.Context) {
		cfg, err := config.LoadConfig()
		if err != nil {
			return
		}
		db, err := database.Connect(c.Request.Context(), cfg)
		if err != nil {
			return
		}
		repo := repositories.NewLocationRepo(db)
		res, err := repo.List(c, 1)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
