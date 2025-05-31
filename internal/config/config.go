package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Env      string   `yaml:"env" env-default:"local"`
	GRPC     GRPC     `yaml:"grpc" env-required:"true"`
	Postgres Postgres `yaml:"postgres" env-required:"true"`
}

type GRPC struct {
	Port    string        `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type Postgres struct {
	Port     string `yaml:"port" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
}

func ParseConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("./configs/local.yaml", &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
