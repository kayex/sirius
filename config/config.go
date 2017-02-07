package config

type AppConfig struct {
	DB   DBConfig
	Maps MapsConfig
}

type MapsConfig struct {
	APIKey string
}

type DBConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}
