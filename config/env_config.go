package config

import (
	"github.com/kayex/sirius/mqtt"
	"os"
)

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
			Host:  os.Getenv("REMOTE_HOST"),
			Token: os.Getenv("REMOTE_TOKEN"),
		},
		Maps: MapsConfig{
			APIKey: os.Getenv("MAPS_API_KEY"),
		},
		MQTT: MQTTConfig{
			Config: mqtt.Config{
				Host: os.Getenv("SYNC_HOST"),
				Port: os.Getenv("SYNC_PORT"),
				User: os.Getenv("SYNC_USER"),
				Pass: os.Getenv("SYNC_PASS"),
				CID:  os.Getenv("SYNC_CID"),
			},
			Topic: os.Getenv("SYNC_TOPIC"),
		},
	}
}
