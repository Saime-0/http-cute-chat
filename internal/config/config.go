package config

import (
	"github.com/BurntSushi/toml"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"os"
	"strings"
)

type Config struct {
	Database             Database `toml:"database"`
	SMTP                 SMTP     `toml:"smtp"`
	Logger               Logger   `toml:"logger"`
	SecretKey            string   `toml:"secret_key"`
	AppPort              string   `toml:"app_port"`
	PasswordSalt         string   `toml:"password_salt"`
	QueryComplexityLimit int      `toml:"query_complexity_limit"`
}

type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DbName   string `toml:"db_name"`
}

type SMTP struct {
	Author string `toml:"author"`
	From   string `toml:"from"`
	Passwd string `toml:"passwd"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
}

type Logger struct {
	Output          clog.Output `toml:"output"`
	Level           string      `toml:"level"`
	MongoDBPassword string      `toml:"mongo_db_password"`
	MongoDBUser     string      `toml:"mongo_db_user"`
	MongoDBCluster  string      `toml:"mongo_db_cluster"`
	DBName          string      `toml:"db_name"`
}

var defaultConfig = &Config{
	Database: Database{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "7050",
		DbName:   "chat_db",
	},
	SMTP: SMTP{
		Author: "http-cute-chat",
		//From:   " ",
		//Passwd: " ",
		Host: "smtp.yandex.ru",
		Port: 0,
	},
	Logger: Logger{
		Output:          clog.MongoDB,
		Level:           "debug",
		MongoDBPassword: os.Getenv("LOGDB_PASSWORD"),
		MongoDBUser:     os.Getenv("LOGDB_USER"),
		MongoDBCluster:  "log-db.5qcx4.mongodb.net",
		DBName:          "logs",
	},
	SecretKey:            os.Getenv("SECRET_SIGNING_KEY"),
	PasswordSalt:         os.Getenv("GLOBAL_PASSWORD_SALT"),
	AppPort:              "8080",
	QueryComplexityLimit: 15,
}

func NewConfig(path string) *Config {
	cfg := &Config{
		Database: Database{},
	}
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return defaultConfig
	}
	if cfg.SecretKey == "" {
		cfg.SecretKey = os.Getenv("SECRET_SIGNING_KEY")
	}
	if cfg.PasswordSalt == "" {
		cfg.PasswordSalt = os.Getenv("GLOBAL_PASSWORD_SALT")
	}
	if cfg.SMTP.From == "" {
		cfg.SMTP.From = os.Getenv("SMTP_EMAIL_LOGIN")
	}
	if cfg.SMTP.Passwd == "" {
		cfg.SMTP.Passwd = os.Getenv("SMTP_EMAIL_PASSWD")
	}
	if cfg.Logger.MongoDBPassword == "" {
		cfg.Logger.MongoDBPassword = os.Getenv("LOGDB_PASSWORD")
	}
	if cfg.Logger.MongoDBUser == "" {
		cfg.Logger.MongoDBUser = os.Getenv("LOGDB_USER")
	}
	cfg.Logger.Level = strings.ToLower(cfg.Logger.Level)

	return cfg
}
