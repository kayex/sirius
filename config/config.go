package config

type Db struct {
	DbUser     string
	DbPassword string
}

type Config struct {
	Db
}
