package handlers

import (
	"field_archive/server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecordingHandler struct {
	service *services.RecordingService
}

func NewRecordingHandler(s *services.RecordingService) *RecordingHandler {
	return &RecordingHandler{service: s}
}

func (h *RecordingHandler) GetByID(c *gin.Context, id int) {
	record, err := h.service.GetByID(id, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch recording"})
		return
	}
	c.JSON(http.StatusOK, record)
}
