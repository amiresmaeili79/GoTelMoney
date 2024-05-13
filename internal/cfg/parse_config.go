package cfg

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// ParseConfig loads configs from .env file
func ParseConfig(path string) Config {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	cfg := Config{}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	return cfg
}
