package postgres

import (
	"database/sql"
	"fmt"

	"github.com/antoniomiletta/fileman/config"
	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func Connect(cfg config.DatabaseConfig) (*DB, error) {
	conn, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	conn.SetMaxOpenConns(cfg.MaxOpenConns)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
