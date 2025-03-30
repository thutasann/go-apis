package config

import "github.com/thutasann/rssagg/internal/database"

// API Config Struct
type APIConfig struct {
	DB *database.Queries
}
