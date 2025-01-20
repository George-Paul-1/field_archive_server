package repositories

import (
	"context"
	"field_archive/server/entities"
	"field_archive/server/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type RecordingRepository interface {
	Insert(recording entities.Recording, ctx context.Context) (int, error)
	GetRowByID(id int, ctx context.Context) (entities.Recording, error)
	Update(recording entities.Recording, ctx context.Context) error
	Delete(id int, ctx context.Context) error
	List(ctx context.Context, limit int) ([]entities.Recording, error)
}

type RecordingRepoImplement struct {
	conn database.Database
}

func NewRecordingRepo(db *database.Postgres) *RecordingRepoImplement {
	return &RecordingRepoImplement{conn: db}
}

func (r *RecordingRepoImplement) Insert(recording entities.Recording, ctx context.Context) (int, error) {

	query := `INSERT INTO recordings` +
		`(title, audio_location, artwork_location, date_uploaded, recording_date, location_id, user_id, ` +
		`duration, format, description, equipment, file_size, channels, license) ` +
		`VALUES ` +
		`(@title, @audio_location, @date_uploaded, @recording_date, @location_id, @user_id, @duration, ` +
		`@format, @description, @equipment, @file_Size, @channels, @license) ` +
		`RETURNING id`
	args := pgx.NamedArgs{
		"title":          recording.Title,
		"audio_location": recording.AudioLocation,
		"date_uploaded":  recording.DateUploaded,
		"recording_date": recording.RecordingDate,
		"location_id":    recording.LocationID,
		"user_id":        recording.UserID,
		"duration":       recording.Duration,
		"format":         recording.Format,
		"description":    recording.Description,
		"equipment":      recording.Equipment,
		"file_size":      recording.Size,
		"channels":       recording.Channels,
		"license":        recording.License,
	}
	var id int
	err := r.conn.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("unable to insert row: %w", err)
	}
	return id, nil
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
		&recording.UserID,
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
	query := `UPDATE recordings ` +
		`SET title = @title, audio_location = @audio_location, ` +
		`artwork_location = @artwork_location, date_uploaded = @date_uploaded, ` +
		`recording_date = @recording_date, location_id = @location_id, user_id = @user_id, duration = @duration, ` +
		`format = @format, description = @description, equipment = @equipment, file_size = @file_size, ` +
		`channels = @channels, license = @license ` +
		`WHERE id = @id`
	args := pgx.NamedArgs{
		"title":          recording.Title,
		"audio_location": recording.AudioLocation,
		"date_uploaded":  recording.DateUploaded,
		"recording_date": recording.RecordingDate,
		"location_id":    recording.LocationID,
		"user_id":        recording.UserID,
		"duration":       recording.Duration,
		"format":         recording.Format,
		"description":    recording.Description,
		"equipment":      recording.Equipment,
		"file_size":      recording.Size,
		"channels":       recording.Channels,
		"license":        recording.License,
		"id":             recording.ID,
	}
	_, err := r.conn.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (r *RecordingRepoImplement) Delete(id int, ctx context.Context) error {
	query := `DELETE FROM recordings WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := r.conn.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("unable to Delete Row: %v", err)
	}
	return nil
}

func (r *RecordingRepoImplement) List(ctx context.Context, limit int) ([]entities.Recording, error) {
	res := []entities.Recording{}
	query := `SELECT * FROM recordings LIMIT $1::int`
	rows, err := r.conn.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		recording := entities.Recording{}
		err := rows.Scan(
			&recording.ID,
			&recording.Title,
			&recording.AudioLocation,
			&recording.ArtworkLocation,
			&recording.DateUploaded,
			&recording.RecordingDate,
			&recording.LocationID,
			&recording.UserID,
			&recording.Duration,
			&recording.Format,
			&recording.Description,
			&recording.Equipment,
			&recording.Size,
			&recording.Channels,
			&recording.License,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, recording)
	}

	return res, nil
}
