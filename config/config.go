package config

type AppConfig struct {
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUser     string
	DBPassword string
	Maps       MapsConfig
}

type MapsConfig struct {
	APIKey string
}
