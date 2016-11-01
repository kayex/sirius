package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JsonConfig struct {
	Config
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
}

func NewJsonConfig(filename string) *Config {
	dat, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Printf("Failed reading JSON config %v\n", filename)
		panic(err)
	}

	var cfg JsonConfig = new(JsonConfig)
	decode(dat, cfg)

	return &cfg
}

func decode(b []byte, cfg JsonConfig) {
	if err := json.Unmarshal(b, &cfg); err != nil {
		panic(err)
	}
}
