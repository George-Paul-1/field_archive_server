package database

import (
	"context"
	"field_archive/server/internal/config"
	"log"

	"github.com/jackc/pgx/v5"
)

func Connect(cfg *config.Config) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), cfg.DB_Url)
	if err != nil {
		log.Fatalf("Error establishing db connection %v", err)
	}
	defer conn.Close(context.Background())
	return conn
}
