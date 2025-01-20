package repositories

import (
	"context"
	"field_archive/server/entities"
	"field_archive/server/internal/database"
	"fmt"
)

// TODO : Build location repository with CRUD functions
type LocationRepository interface {
	Insert(recording entities.Location, ctx context.Context) (int, error)
	GetRowByID(id int, ctx context.Context) (entities.Location, error)
	Update(recording entities.Location, ctx context.Context) error
	Delete(id int, ctx context.Context) error
	List(ctx context.Context, limit int) ([]entities.Location, error)
}

type LocationRepoImplement struct {
	conn database.Database
}

func NewLocationRepo(db *database.Postgres) *LocationRepoImplement {
	return &LocationRepoImplement{conn: db}
}

func (r *LocationRepoImplement) Insert(location entities.Location, ctx context.Context) (int, error) {
	query := `INSERT INTO locations ` +
		`(name, description, geom) ` +
		`VALUES (@name, @description, ST_SetSRID(ST_MakePoint(@longitude, @latitude), 4326)) ` +
		`RETURNING id`
	args := map[string]interface{}{
		"name":        location.Name,
		"description": location.Description,
		"longitude":   location.Longitude,
		"latitude":    location.Latitude,
	}
	var id int
	err := r.conn.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("unable to insert row: %w", err)
	}
	return id, nil
}
