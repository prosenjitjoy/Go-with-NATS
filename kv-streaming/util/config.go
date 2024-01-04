package util

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDatabase struct {
	Username string        `env:"NATS_USER" env-required:"true"`
	Password string        `env:"NATS_PASS" env-required:"true"`
	HostPort string        `env:"HOST_PORT" env-required:"true"`
	NatsUrl  string        `env:"NATS_URL" env-required:"true"`
	Timeout  time.Duration `env:"TIMEOUT" env-required:"true"`
}

func LoadConfig(path string) (*ConfigDatabase, error) {
	var cfg ConfigDatabase

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return &ConfigDatabase{}, err
	}

	return &cfg, nil
}
