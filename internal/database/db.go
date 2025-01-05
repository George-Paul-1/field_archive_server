package database

import (
	"context"
	"field_archive/server/internal/config"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func Connect(ctx context.Context, cfg *config.Config) (*Postgres, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, cfg.DB_Url)
		if err != nil {
			log.Fatalf("Error establishing db connection %v", err)
		}
		pgInstance = &Postgres{db}
	})
	return pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.DB.Close()
}
