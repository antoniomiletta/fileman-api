package postgres

import (
	"fmt"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/domain"
)

type DB struct {
	conn *gorm.DB
}

func Connect(cfg config.DatabaseConfig) (*DB, error) {
	db, err := gorm.Open(pg.Open(cfg.URL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	db.AutoMigrate(domain.Models...)

	conn, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get conn: %w", err)
	}

	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &DB{conn: db}, nil
}

func (db *DB) Close() error {
	conn, err := db.conn.DB()
	if err != nil {
		return err
	}
	return conn.Close()
}
