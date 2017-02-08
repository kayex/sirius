package config

import "os"

func FromEnv() AppConfig {
	return AppConfig{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Database: os.Getenv("DB_DATABASE"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Remote: RemoteConfig{
			URL:   os.Getenv("REMOTE_URL"),
			Token: os.Getenv("REMOTE_TOKEN"),
		},
		Maps: MapsConfig{
			APIKey: os.Getenv("MAPS_API_KEY"),
		},
	}
}
