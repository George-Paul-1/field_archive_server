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

type MockRow struct {
	mockScan func(dest ...any) error
}

func (r *MockRow) Scan(dest ...any) error {
	return r.mockScan(dest)
}

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

func TestGetRowByID(t *testing.T) {
	check := `SELECT * FROM recordings WHERE id=@id`
	expectedRecording := entities.Recording{
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
	mockDB := &MockDatabase{
		mockQueryRow: func(ctx context.Context, query string, args ...any) pgx.Row {
			if check == query {
				return &MockRow{mockScan: func(dest ...any) error {
					innerSlice, ok := dest[0].([]any)
					if !ok {
						return errors.New("Can't access inner slice")
					}
					*(innerSlice[0].(*int)) = 1
					*(innerSlice[1].(*string)) = "Test Title"
					*(innerSlice[2].(*string)) = "test/audio/location.mp3"
					*(innerSlice[3].(**string)) = nil
					*(innerSlice[4].(**time.Time)) = nil
					*(innerSlice[5].(*time.Time)) = time.Date(2025, 1, 6, 20, 2, 57, 0, time.UTC)
					*(innerSlice[6].(*int)) = 1
					*(innerSlice[7].(*int)) = 120
					*(innerSlice[8].(*string)) = "mp3"
					*(innerSlice[9].(*string)) = "test description"
					*(innerSlice[10].(*string)) = "test equipment"
					*(innerSlice[11].(*float64)) = 2048
					*(innerSlice[12].(*string)) = "2"
					*(innerSlice[13].(*string)) = "Creative Commons"
					return nil
				}}
			}
			return &MockRow{mockScan: func(dest ...any) error {
				return errors.New("row not found")
			}}
		},
	}
	repo := &RecordingRepoImplement{conn: mockDB}
	test, err := repo.GetRowByID(1, context.Background())
	if test != expectedRecording {
		t.Errorf("GetRowByID failed: return = %v, but expected %v", test, expectedRecording)
	}
	if err != nil {
		t.Errorf("getRowByID failed: Error returned from func %v", err)
	}

}
