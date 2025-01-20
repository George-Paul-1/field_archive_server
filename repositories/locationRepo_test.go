package repositories

import (
	"context"
	"errors"
	"field_archive/server/entities"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestInsertLocation(t *testing.T) {
	check := `INSERT INTO locations ` +
		`(name, description, geom) ` +
		`VALUES (@name, @description, ST_SetSRID(ST_MakePoint(@longitude, @latitude), 4326)) ` +
		`RETURNING id`

	mockDB := &MockDatabase{
		mockQueryRow: func(ctx context.Context, query string, args ...any) pgx.Row {
			if check == query {
				return &MockRow{mockScan: func(dest ...any) error {
					innerSlice, ok := dest[0].([]any)
					if !ok {
						return errors.New("Unable to access inner slice")
					}
					*(innerSlice[0].(*int)) = 3
					return nil
				}}
			}
			return &MockRow{mockScan: func(dest ...any) error {
				return errors.New("Row not found")
			},
			}
		},
	}
	repo := &LocationRepoImplement{conn: mockDB}
	longitude := "1.00"
	latitude := "2.00"
	location := entities.Location{
		Name:        "Test Location",
		Description: "Test Description",
		Longitude:   &longitude,
		Latitude:    &latitude,
	}
	id, err := repo.Insert(location, context.Background())
	assert.NoError(t, err)
	if id != 3 {
		t.Fatalf("INSERT ERROR: returning ids not equal")
	}
}

func TestGetLocationByID(t *testing.T) {
	check := `SELECT id, name, description, ST_AsGeoJSON(geom) AS geom FROM locations WHERE id = @id`
	mockDB := &MockDatabase{
		mockQueryRow: func(ctx context.Context, query string, args ...any) pgx.Row {
			if check == query {
				return &MockRow{mockScan: func(dest ...any) error {
					innerSlice, ok := dest[0].([]any)
					if !ok {
						return errors.New("Unable to access inner slice")
					}
					*(innerSlice[0].(*int)) = 1
					*(innerSlice[1].(*string)) = "Test Location"
					*(innerSlice[2].(*string)) = "Test Description"
					*(innerSlice[3].(*string)) = "Test Geom"
					return nil
				}}
			}
			return &MockRow{mockScan: func(dest ...any) error {
				return errors.New("Row not found")
			},
			}
		},
	}
	repo := &LocationRepoImplement{conn: mockDB}
	location, err := repo.GetRowByID(1, context.Background())
	assert.NoError(t, err)
	if location.ID != 1 {
		t.Fatalf("GET ERROR: returning ids not equal")
	}
}

func TestUpdateLocation(t *testing.T) {
	check := `UPDATE locations SET name = @name, description = @description, geom = ST_SetSRID(ST_MakePoint(@longitude, @latitude), 4326) WHERE id = @id`
	mockDB := &MockDatabase{
		mockExec: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
			if check == query {
				return pgconn.CommandTag{}, nil
			}
			return pgconn.CommandTag{}, errors.New("Row not found")
		},
	}
	repo := &LocationRepoImplement{conn: mockDB}
	longitude := "1.00"
	latitude := "2.00"
	location := entities.Location{
		ID:          1,
		Name:        "Test Location",
		Description: "Test Description",
		Longitude:   &longitude,
		Latitude:    &latitude,
	}
	err := repo.Update(location, context.Background())
	assert.NoError(t, err)
}

func TestDeleteLocation(t *testing.T) {
	check := `DELETE FROM locations WHERE id = @id`
	mockDB := MockDatabase{mockExec: func(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
		if check == query {
			return pgconn.CommandTag{}, nil
		} else {
			return pgconn.CommandTag{}, errors.New("DELETE ERROR: Query did not match check")
		}
	}}
	repo := &LocationRepoImplement{conn: &mockDB}
	err := repo.Delete(1, context.Background())
	assert.NoError(t, err)
}
