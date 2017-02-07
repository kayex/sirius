package config

import "os"

func FromEnv() AppConfig {
	return AppConfig{
		DB: DBConfig{
			Host:     env(`DB_HOST`, `127.0.0.1`),
			Port:     env(`DB_PORT`, `5432`),
			Database: env(`DB_DATABASE`, ``),
			User:     env(`DB_USER`, ``),
			Password: env(`DB_PASSWORD`, ``),
		},
		Maps: MapsConfig{
			APIKey: env("MAPS_API_KEY", ""),
		},
	}
}

func env(key string, def string) string {
	val := os.Getenv(key)

	if val == "" {
		val = def
	}

	return val
}
