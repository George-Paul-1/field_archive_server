package repositories

import (
	"context"
	"field_archive/server/entities"
	"field_archive/server/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type RecordingRepository interface {
	Insert(recording entities.Recording, ctx context.Context) error
	GetRowByID(id int, ctx context.Context) (entities.Recording, error)
	Update(recording entities.Recording, ctx context.Context) error
	Delete(id int, ctx context.Context) error
	List(ctx context.Context) ([]entities.Recording, error)
}

type RecordingRepoImplement struct {
	conn database.Database
}

func NewRecordingRepo(db *database.Postgres) *RecordingRepoImplement {
	return &RecordingRepoImplement{conn: db}
}

func (r *RecordingRepoImplement) Insert(recording entities.Recording, ctx context.Context) error {

	query := `INSERT INTO recordings` +
		`(id, title, audio_location, artwork_location, date_uploaded, recording_date, location_id, ` +
		`duration, format, description, equipment, file_size, channels, license) ` +
		`VALUES ` +
		`(@title, @audio_location, @date_uploaded, @recording_date, @location_id, @duration, @format, @description, @equipment, @file_Size, @channels, @license)`
	args := map[string]interface{}{
		"title":          recording.Title,
		"audio_location": recording.AudioLocation,
		"date_uploaded":  recording.DateUploaded,
		"recording_date": recording.RecordingDate,
		"location_id":    recording.LocationID,
		"duration":       recording.Duration,
		"format":         recording.Format,
		"description":    recording.Description,
		"equipment":      recording.Equipment,
		"file_size":      recording.Size,
		"channels":       recording.Channels,
		"license":        recording.License,
	}
	_, err := r.conn.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (r *RecordingRepoImplement) GetRowByID(id int, ctx context.Context) (entities.Recording, error) {
	query := `SELECT * FROM recordings WHERE id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}
	var recording entities.Recording
	err := r.conn.QueryRow(ctx, query, args).Scan(
		&recording.ID,
		&recording.Title,
		&recording.AudioLocation,
		&recording.ArtworkLocation,
		&recording.DateUploaded,
		&recording.RecordingDate,
		&recording.LocationID,
		&recording.Duration,
		&recording.Format,
		&recording.Description,
		&recording.Equipment,
		&recording.Size,
		&recording.Channels,
		&recording.License,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No rows found for the given ID
			return entities.Recording{}, fmt.Errorf("recording with id %d not found", id)
		}
		return entities.Recording{}, fmt.Errorf("unable to fetch recording %w", err)
	}
	return recording, nil
}

func (r *RecordingRepoImplement) Update(recording entities.Recording, ctx context.Context) error {
	return nil
}

func (r *RecordingRepoImplement) Delete(id int, ctx context.Context) error {
	return nil
}

func (r *RecordingRepoImplement) List(ctx context.Context) ([]entities.Recording, error) {
	slice := []entities.Recording{}
	return slice, nil
}
