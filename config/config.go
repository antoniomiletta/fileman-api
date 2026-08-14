package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server  ServerConfig
	DB      DBConfig
	Storage StorageConfig
	Auth    AuthConfig
	Cleanup CleanupConfig
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
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
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
	Backend     string `env:"STORAGE_BACKEND" envDefault:"local"`
	LocalConfig LocalConfig
	S3Config    S3Config
}

type LocalConfig struct {
	LocalRoot string `env:"LOCAL_STORAGE_ROOT"`
}

type S3Config struct {
	S3Bucket string `env:"S3_BUCKET"`
}

type AuthConfig struct {
	JWTSecret   string        `env:"JWT_SECRET"`
	TokenExpiry time.Duration `env:"TOKEN_EXPIRY" envDefault:"24h"`
}

type CleanupConfig struct {
	PollInterval    time.Duration `env:"CLEANUP_POLL_INTERVAL" envDefault:"10s"`
	BatchSize       int           `env:"CLEANUP_BATCH_SIZE" envDefault:"20"`
	MaxInFlight     int           `env:"CLEANUP_MAX_IN_FLIGHT" envDefault:"10"`
	JobTimeout      time.Duration `env:"CLEANUP_JOB_TIMEOUT" envDefault:"30s"`
	MaxAttempts     int           `env:"CLEANUP_MAX_ATTEMPTS" envDefault:"5"`
	ReclaimInterval time.Duration `env:"CLEANUP_RECLAIM_INTERVAL" envDefault:"5m"`
	StaleAfter      time.Duration `env:"CLEANUP_STALE_AFTER" envDefault:"10m"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
