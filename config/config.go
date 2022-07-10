package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	ConfigEnv     = "CHATTY_CONFIG"
	DefaultConfig = "./conf/config-dev.yml"
)

type Database struct {
	Password string `mapstructure:"password"`
}

type Log struct {
	Level string `mapstructure:"level"`
}

type Server struct {
	Port int `mapstructure:"port"`
	Log  Log `mapstructure:"log"`
}

type Auth struct {
	Issuer      string `mapstructure:"issuer"`
	ClientID    string `mapstructure:"client_id"`
	SecretID    string `mapstructure:"secret_id"`
	RedirectURL string `mapstructure:"redirect_url"`
	State       string `mapstructure:"state"`
}

type Config struct {
	Server   Server   `mapstructure:"server"`
	Auth     Auth     `mapstructure:"auth"`
	Database Database `mapstructure:"database"`
}

func Load() (*Config, error) {
	confPath := DefaultConfig
	confEnvPath := os.Getenv(ConfigEnv)
	if confEnvPath != "" {
		confPath = confEnvPath
	}

	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to load config, error: %v", err)
	}

	viper.SetConfigFile(confPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}