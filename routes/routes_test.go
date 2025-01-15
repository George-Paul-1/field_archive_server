package routes

import (
	"context"
	"field_archive/server/entities"
	"field_archive/server/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	mockGetByID   func(id int) (entities.Recording, error)
	mockListItems func(limit int, ctx context.Context) ([]entities.Recording, error)
}

func (m *mockService) GetByID(id int, ctx context.Context) (entities.Recording, error) {
	return m.mockGetByID(id)
}

func (m *mockService) ListItems(limit int, ctx context.Context) ([]entities.Recording, error) {
	return m.mockListItems(limit, ctx)
}

func TestTestRoute(t *testing.T) {
	router := gin.Default()

	h := handlers.RecordingHandler{}
	DefineRoutes(router, &h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message": "test"}`, w.Body.String())
}

func TestRecordingsGetByIDRoute(t *testing.T) {
	router := gin.Default()

	mockResponse := entities.Recording{
		ID:              1,
		Title:           "Test Title",
		AudioLocation:   "test/audio/location.mp3",
		ArtworkLocation: nil,
		DateUploaded:    nil,
		RecordingDate:   time.Date(2025, 1, 6, 20, 2, 57, 0, time.UTC),
		LocationID:      1,
		Duration:        120,
		Format:          "mp3",
		Description:     "test description",
		Equipment:       "test equipment",
		Size:            2048,
		Channels:        "2",
		License:         "Creative Commons",
	}
	mockService := &mockService{
		mockGetByID: func(id int) (entities.Recording, error) {
			return mockResponse, nil
		},
	}
	h := handlers.RecordingHandler{Service: mockService}
	DefineRoutes(router, &h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/recordings/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{
  "ID": 1,
  "Title": "Test Title",
  "AudioLocation": "test/audio/location.mp3",
  "ArtworkLocation": null,
  "DateUploaded": null,
  "RecordingDate": "2025-01-06T20:02:57Z",
  "LocationID": 1,
  "Duration": 120,
  "Format": "mp3",
  "Description": "test description",
  "Equipment": "test equipment",
  "Size": 2048,
  "Channels": "2",
  "License": "Creative Commons"
}`, w.Body.String())
}

func TestListGETRoute(t *testing.T) {
	router := gin.Default()

	mockResponse := []entities.Recording{
		{
			ID:              1,
			Title:           "Test Title",
			AudioLocation:   "test/audio/location.mp3",
			ArtworkLocation: nil,
			DateUploaded:    nil,
			RecordingDate:   time.Date(2025, 1, 6, 20, 2, 57, 0, time.UTC),
			LocationID:      1,
			Duration:        120,
			Format:          "mp3",
			Description:     "test description",
			Equipment:       "test equipment",
			Size:            2048,
			Channels:        "2",
			License:         "Creative Commons",
		},
	}
	mockService := &mockService{
		mockListItems: func(limit int, ctx context.Context) ([]entities.Recording, error) {
			return mockResponse, nil
		},
	}
	h := handlers.RecordingHandler{Service: mockService}
	DefineRoutes(router, &h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/recordings/list/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `[{
  "ID": 1,
  "Title": "Test Title",
  "AudioLocation": "test/audio/location.mp3",
  "ArtworkLocation": null,
  "DateUploaded": null,
  "RecordingDate": "2025-01-06T20:02:57Z",
  "LocationID": 1,
  "Duration": 120,
  "Format": "mp3",
  "Description": "test description",
  "Equipment": "test equipment",
  "Size": 2048,
  "Channels": "2",
  "License": "Creative Commons"
}]`, w.Body.String())
}
