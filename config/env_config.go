package config

import "os"

func FromEnv() Config {
	return Config{
		DbHost:     env(`DB_HOST`, `127.0.0.1`),
		DbPort:     env(`DB_PORT`, `5432`),
		DbDatabase: env(`DB_DATABASE`, ``),
		DbUser:     env(`DB_USER`, ``),
		DbPassword: env(`DB_PASSWORD`, ``),
	}
}

func env(key string, def string) string {
	val := os.Getenv(key)

	if val == "" {
		val = def
	}

	return val
}
