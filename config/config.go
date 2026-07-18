package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string `env:"PORT"`
}

type DatabaseConfig struct {
	URL             string        `env:"PORT"`
	MaxOpenConns    int           `env:"DB_URL"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME"`
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
