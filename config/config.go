package config

type AppConfig struct {
	DB     DBConfig
	Remote RemoteConfig
	Maps   MapsConfig
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

type RemoteConfig struct {
	URL   string
	Token string
}
