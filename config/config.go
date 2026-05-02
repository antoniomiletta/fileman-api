package config

import "time"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL          string
	MaxOpenConns int
}

type StorageConfig struct {
	Backend          string
	LocalStoragePath string
}

type AuthConfig struct {
	JWTSecret   string
	TokenExpiry time.Duration
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			URL:          mustGetEnv("DATABASE_URL"),
			MaxOpenConns: getInt("DB_MAX_OPEN_CONNS", 20),
		},
		Storage: StorageConfig{
			Backend:          getEnv("STORAGE_BACKEND", "local"),
			LocalStoragePath: getEnv("LOCAL_STORAGE_PATH", "./data"),
		},
		Auth: AuthConfig{
			JWTSecret:   mustGetEnv("JWT_SECRET"),
			TokenExpiry: getDuration("TOKEN_EXPIRY", 24*time.Hour),
		},
	}
}
