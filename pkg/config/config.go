package config

import (
	"lila_games/internal/models"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Get reads apps configuration by reading ENV vars
func Get() (*models.Config, error) {
	if err := godotenv.Overload(); err != nil {
		return nil, err
	}
	var cfg models.Config
	if err := envconfig.Process("LILA", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
