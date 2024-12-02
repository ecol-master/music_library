package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	TokenTTL  time.Duration `yaml:"token_ttl" env-required:"true"`
	Secret    string        `yaml:"secret" env-required:"true"`
	App       AppConfig
	APIServer APIServerConfig
	APIClient APIClientConfig
	Postgres  PostgresConfig
}

type AppConfig struct {
	Port uint16 `yaml:"port" env:"PORT" env-default:"8080"`
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
}

type PostgresConfig struct {
	Port           string        `yaml:"port" env:"PORT" env-default:"5432"`
	Host           string        `yaml:"host" env:"HOST" env-default:"localhost"`
	Name           string        `yaml:"name" env:"NAME" env-default:"music_lib"`
	User           string        `yaml:"user" env:"USER" env-default:"user"`
	Password       string        `yaml:"password" env:"PASSWORD"`
	ConnTimeExceed time.Duration `yaml:"conn_time_exceed" env:"CONN_TIME_EXCEED" env-default:"3s"`
}

type APIClientConfig struct {
	BaseURL        string        `yaml:"base_url" env:"BASE_URL" env-required:"true"`
	RequestTimeout time.Duration `yaml:"request_timeout" env:"REQUEST_TIMEOUT" env-default:"5s"`
}

type APIServerConfig struct {
}

// Function will panic if can not read config file or environment variables
func MustLoad() *Config {
	var cfg Config

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(errors.Wrap(err, "config file dost not exists"))
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(errors.Wrap(err, "failed to read config"))
	}
	return &cfg
}

// fetchConfigPath fetches config path from command line flar or env variable
// Priority: flag > env > default
// Default value is empty string
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path fo config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
