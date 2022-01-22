package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Database     Database `toml:"database"`
	SecretKey    string   `toml:"secret_key"`
	AppPort      string   `toml:"app_port"`
	PasswordSalt string   `toml:"password_salt"`
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
	SecretKey:    os.Getenv("SECRET_SIGNING_KEY"),
	PasswordSalt: os.Getenv("GLOBAL_PASSWORD_SALT"),
	AppPort:      "8080",
}

func NewConfig(path string) *Config {
	cfg := &Config{
		Database: Database{},
	}
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		println("Default config loaded")
		return defaultConfig
	}
	if cfg.SecretKey == "" {
		cfg.SecretKey = os.Getenv("SECRET_SIGNING_KEY")
	}
	if cfg.PasswordSalt == "" {
		cfg.PasswordSalt = os.Getenv("GLOBAL_PASSWORD_SALT")
	}
	println(cfg.SecretKey, cfg.PasswordSalt)
	println("Configure file found and loaded")
	return cfg
}
