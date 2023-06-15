package configs

type Config struct {
	Server   Server   `yaml:"Server"`
	Postgres Postgres `yaml:"Postgres"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Postgres struct {
	DB             string `yaml:"db"`
	ConnectionToDB string `yaml:"connectionToDB"`
}
