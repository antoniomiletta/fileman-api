package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server   ServerConfig
	Database DBConfig
	Storage  StorageConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string `env:"PORT"`
}

type DBConfig struct {
	User            string        `env:"DB_USER"`
	Password        string        `env:"DB_PASSWORD"`
	Host            string        `env:"DB_HOST"`
	Port            int           `env:"DB_PORT"`
	Name            string        `env:"DB_NAME"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME"`
}

func (c *DBConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Name,
	)
}

type StorageConfig struct {
	Backend  string `env:"STORAGE_BACKEND"`
	S3Bucket string `env:"S3_BUCKET"`
}

type AuthConfig struct {
	JWTSecret   string        `env:"JWT_SECRET"`
	TokenExpiry time.Duration `env:"TOKEN_EXPIRY"`
}

func Load() *Config {
	var cfg Config
	env.Parse(&cfg)

	return &cfg
}
