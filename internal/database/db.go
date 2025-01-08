package database

import (
	"context"
	"field_archive/server/internal/config"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

type Postgres struct {
	DB *pgxpool.Pool
}

func (p *Postgres) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return p.DB.Exec(ctx, query, args...)
}

func (p *Postgres) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return p.DB.QueryRow(ctx, query, args...)
}

func (p *Postgres) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return p.DB.Query(ctx, query, args)
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
