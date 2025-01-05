package repositories

// import (
// 	"context"
// 	"field_archive/server/entities"
// 	"testing"
// 	"time"

// 	"github.com/pashagolub/pgxmock"
// )

// func TestInsert(t *testing.T) {
// 	mock, err := pgxmock.NewConn()
// 	if err != nil {
// 		t.Fatalf("failed to create pgxMock", err)
// 	}
// 	defer mock.Close(context.Background())

// 	recording := entities.Recording{
// 		Title:         "Test Title",
// 		AudioLocation: "test/audio/location.mp3",
// 		DateUploaded:  time.Now(),
// 		RecordingDate: time.Now(),
// 		LocationID:    1,
// 		Duration:      120,
// 		Format:        "mp3",
// 		Description:   "Test description",
// 		Equipment:     "Test equipment",
// 		Size:          2048,
// 		Channels:      2,
// 		License:       "Creative Commons",
// 	}

// 	query := `INSERT INTO recordings` +
// 	`\(id, title, audio_location, artwork_location, date_uploaded, recording_date, location_id,` +
// 	` duration, format, description, equipment, file_size, channels, license\) ` +
// 	`VALUES` +
// 	`\(@title, @audio_location, @date_uploaded, @recording_date, @location_id, @duration, @format, @description, @equipment, @file_Size, @channels, @license\)`
// 	mock.ExpectExec(query).
// 		WithArgs(
// 			recording.Title,
// 			recording.AudioLocation,
// 			recording.DateUploaded,
// 			recording.RecordingDate,
// 			recording.LocationID,
// 			recording.Duration,
// 			recording.Format,
// 			recording.Description,
// 			recording.Equipment,
// 			recording.Size,
// 			recording.Channels,
// 			recording.License,
// 		).
// 	WillReturnResult(pgxmock.NewResult("INSERT", 1))
// 	repo := &RecordingRepoImplement{conn: mock,}
// }
