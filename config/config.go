package config

import "github.com/kayex/sirius/mqtt"

type AppConfig struct {
	DB     DBConfig
	Remote RemoteConfig
	Maps   MapsConfig
	MQTT   MQTTConfig
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

type MQTTConfig struct {
	mqtt.Config
	Topic string
}
