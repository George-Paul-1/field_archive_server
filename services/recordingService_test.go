package services

import (
	"context"
	"field_archive/server/entities"
	"field_archive/server/internal/database"
	"field_archive/server/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	mockInsert     func(recording entities.Recording, ctx context.Context) (int, error)
	mockGetRowByID func(id int, ctx context.Context) (entities.Recording, error)
	mockUpdate     func(recording entities.Recording, ctx context.Context) error
	mockDelete     func(id int, ctx context.Context) error
	mockList       func(ctx context.Context, limit int) ([]entities.Recording, error)
}

func (r *mockRepo) Insert(recording entities.Recording, ctx context.Context) (int, error) {
	return r.mockInsert(recording, ctx)
}

func (r *mockRepo) GetRowByID(id int, ctx context.Context) (entities.Recording, error) {
	return r.mockGetRowByID(id, ctx)
}

func (r *mockRepo) Update(recording entities.Recording, ctx context.Context) error {
	return r.mockUpdate(recording, ctx)
}

func (r *mockRepo) Delete(id int, ctx context.Context) error {
	return r.mockDelete(id, ctx)
}

func (r *mockRepo) List(ctx context.Context, limit int) ([]entities.Recording, error) {
	return r.mockList(ctx, limit)
}

func TestNewRecordingService(t *testing.T) {
	r := repositories.NewRecordingRepo(&database.Postgres{})
	s := NewRecordingService(r)
	check := &RecordingService{repo: r}
	assert.Equal(t, s, check)
}

func TestGetByID(t *testing.T) {
	check := entities.Recording{
		Title:         "Test Title",
		AudioLocation: "test/audio/location.mp3",
		DateUploaded:  nil,
		RecordingDate: time.Now(),
		LocationID:    1,
		Duration:      120,
		Format:        "mp3",
		Description:   "Test description",
		Equipment:     "Test equipment",
		Size:          2048,
		Channels:      "2",
		License:       "Creative Commons",
	}

	mockRepo := &mockRepo{
		mockGetRowByID: func(id int, ctx context.Context) (entities.Recording, error) {
			return check, nil
		},
	}
	s := &RecordingService{repo: mockRepo}
	res, err := s.GetByID(1, context.Background())
	assert.Equal(t, check, res)
	if err != nil {
		t.Errorf("Error retrieving by ID: %v", err)
	}
}

func TestListItems(t *testing.T) {
	check := []entities.Recording{
		{Title: "Test Title",
			AudioLocation: "test/audio/location.mp3",
			DateUploaded:  nil,
			RecordingDate: time.Now(),
			LocationID:    1,
			Duration:      120,
			Format:        "mp3",
			Description:   "Test description",
			Equipment:     "Test equipment",
			Size:          2048,
			Channels:      "2",
			License:       "Creative Commons"},
	}
	mockRepo := &mockRepo{
		mockList: func(ctx context.Context, limit int) ([]entities.Recording, error) {
			return check, nil
		},
	}
	s := &RecordingService{repo: mockRepo}
	res, err := s.ListItems(1, context.Background())
	assert.Equal(t, res, check)
	if err != nil {
		t.Errorf("Error listing items %v", err)
	}
}
