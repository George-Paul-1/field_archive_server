package repositories

import (
	"context"
	"errors"
	"field_archive/server/entities"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

type MockDatabase struct {
	mockExec     func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	mockQueryRow func(ctx context.Context, query string, args ...any) pgx.Row
}

func (m *MockDatabase) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return m.mockExec(ctx, query, args...)
}
func (m *MockDatabase) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return m.mockQueryRow(ctx, query, args...)
}

func TestInsert(t *testing.T) {
	check := `INSERT INTO recordings` +
		`(id, title, audio_location, artwork_location, date_uploaded, recording_date, location_id, ` +
		`duration, format, description, equipment, file_size, channels, license) ` +
		`VALUES ` +
		`(@title, @audio_location, @date_uploaded, @recording_date, @location_id, @duration, @format, @description, @equipment, @file_Size, @channels, @license)`

	mockDB := &MockDatabase{
		mockExec: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
			if check == query {
				return pgconn.CommandTag{}, nil
			}
			return pgconn.CommandTag{}, errors.New("query mismatch")
		},
	}
	repo := &RecordingRepoImplement{conn: mockDB}
	recording := entities.Recording{
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
	err := repo.Insert(recording, context.Background())
	assert.NoError(t, err)
}
