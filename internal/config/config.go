package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Database  Database `toml:"database"`
	SecretKey string   `toml:"secret_key"`
	AppPort   string   `toml:"app_port"`
}

type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DbName   string `toml:"db_name"`
}

var defaultConfig = &Config{
	Database: Database{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "7050",
		DbName:   "chat_db",
	},
	SecretKey: os.Getenv("SECRET_SIGNING_KEY"),
	AppPort:   "8080",
}

func NewConfig(path string) *Config {
	cfg := &Config{
		Database:  Database{},
		SecretKey: "",
	}
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		println("Default config loaded")
		return defaultConfig
	}
	if cfg.SecretKey == "" {
		cfg.SecretKey = os.Getenv("SECRET_SIGNING_KEY")
	}
	println("Configure file found and loaded")
	return cfg
}
