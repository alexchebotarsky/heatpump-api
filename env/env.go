package env

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	envconfig "github.com/sethvargo/go-envconfig"
)

type Config struct {
	LogLevel  slog.Level `env:"LOG_LEVEL,default=debug"`
	LogFormat string     `env:"LOG_FORMAT,default=json"`

	Host string `env:"HOST,default=0.0.0.0"`
	Port uint16 `env:"PORT,default=8000"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var c Config

	// We are loading env variables from .env file only for local development
	err := godotenv.Load(".env")
	if err != nil {
		slog.Debug(fmt.Sprintf("error loading .env file: %v", err))
	}

	err = envconfig.Process(ctx, &c)
	if err != nil {
		return nil, fmt.Errorf("error processing environment variables: %v", err)
	}

	return &c, nil
}
