package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ServerConfig struct {
	Port            string        `envconfig:"SERVER_PORT" default:"8080"`
	ShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"15s"`
}

type LoggerConfig struct {
	Mode string `envconfig:"LOG_MODE" default:"dev"`
}

type DatabaseConfig struct {
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"postgres"`
	Host     string `envconfig:"DB_HOST" default:"postgres"`
	Port     string `envconfig:"DB_PORT" default:"5435"`
	Name     string `envconfig:"DB_NAME" default:"postgres"`
}

func (c *DatabaseConfig) DSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.Name
}

type PoolConfig struct {
	MaxConns    int32         `envconfig:"POOL_MAX_CONNS" default:"10"`
	MinConns    int32         `envconfig:"POOL_MIN_CONNS" default:"2"`
	MaxIdleTime time.Duration `envconfig:"POOL_MAX_IDLE_TIME" default:"1h"`
	MaxLifeTime time.Duration `envconfig:"POOL_MAX_LIFE_TIME" default:"10m"`
}

type JWTConfig struct {
	SecretKey     string        `envconfig:"JWT_SECRET_KEY" default:"secret"`
	TokenDuration time.Duration `envconfig:"JWT_TOKEN_DURATION" default:"15m"`
}

type Config struct {
	Server   ServerConfig
	Logger   LoggerConfig
	Database DatabaseConfig
	Pool     PoolConfig
	JWT      JWTConfig
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load env: %w", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
