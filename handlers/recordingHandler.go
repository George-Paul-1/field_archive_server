package handlers

import (
	"field_archive/server/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecordingHandler struct {
	service *services.RecordingService
}

func NewRecordingHandler(s *services.RecordingService) *RecordingHandler {
	return &RecordingHandler{service: s}
}

func (h *RecordingHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID must be a valid integer",
		})
		return
	}
	record, err := h.service.GetByID(id, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch recording"})
		return
	}
	c.JSON(http.StatusOK, record)
}

func (h *RecordingHandler) ListItems(c *gin.Context) {
	Param := c.Param("limit")
	limit, err := strconv.Atoi(Param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limit must be valid integer",
		})
		return
	}
	recordings, err := h.service.ListItems(limit, c.Request.Context())
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to retrieve items",
		})
		return
	}
	c.JSON(http.StatusOK, recordings)
}
