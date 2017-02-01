package config

import "os"

func FromEnv() Config {
	return Config{
		DBHost:     env(`DB_HOST`, `127.0.0.1`),
		DBPort:     env(`DB_PORT`, `5432`),
		DBDatabase: env(`DB_DATABASE`, ``),
		DBUser:     env(`DB_USER`, ``),
		DBPassword: env(`DB_PASSWORD`, ``),
	}
}

func env(key string, def string) string {
	val := os.Getenv(key)

	if val == "" {
		val = def
	}

	return val
}
