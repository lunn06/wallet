package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config and inner structs describe yaml config field
type Config struct {
	HTTPServer `yaml:"http_server"`
	Database   `yaml:"database"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"127.0.0.1"`
	Port    string `yaml:"port" env-default:"8080"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
}

// ReadConfig parse config from path
func ReadConfig(configPath string) (Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Config{}, err
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
