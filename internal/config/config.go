package config

import (
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Host        string `env:"HOST" envDefault:"localhost"`
	Port        string `env:"SERVER_PORT" envDefault:"8080"`
	AllowOrigin string `env:"ALLOW_ORIGIN" envDefault:"*"`
	Env         string `env:"ENV" envDefault:"dev"`

	Secret        string        `env:"SECRET" envDefault:"SecretSecretSecret"`
	TokenLifespan time.Duration `env:"TOKEN_LIFESPAN" envDefault:"24h"`
	DatabaseUrl   string        `env:"DATABASE_URL" envDefault:"postgres://admin:password@localhost:5432/internship"`
	RedisUrl      string        `env:"REDIS_URL" envDefault:"redis://@localhost:6379/"`
}

func NewConfig() (Config, error) {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Print(".env file could not be imported")
	}

	err = env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
