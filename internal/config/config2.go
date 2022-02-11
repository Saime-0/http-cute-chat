package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"os"
)

type Config2 struct {
	FromEnv
	FromCfgFile
}

func NewConfig2(pathToCfgFile string) (*Config2, error) {
	fromFile := new(FromCfgFile)
	_, err := toml.DecodeFile(pathToCfgFile, fromFile)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось декодировать файл")
	}

	if !fromFile.validate() {
		return nil, errors.New("в файле конфигурации заполнены не все поля")
	}

	if clog.Exists(clog.LogLevel(*fromFile.Logging.LoggingLevel)) {

	}

	fromEnv := &FromEnv{
		PostgresConnection: os.Getenv("POSTGRES_CONNECTION"),
		GlobalPasswordSalt: os.Getenv("GLOBAL_PASSWORD_SALT"),
		MongoDBUri:         os.Getenv("MONGODB_URI"),
		SecretSigningKey:   os.Getenv("SECRET_SIGNING_KEY"),
		SmtpHost:           os.Getenv("SMTP_HOST"),
		SmtpEmailLogin:     os.Getenv("SMTP_EMAIL_LOGIN"),
		SmtpEmailPasswd:    os.Getenv("SMTP_EMAIL_PASSWD"),
	}

	if !fromEnv.validate() {
		return nil, errors.New("не установлены некоторые переменные окружения")
	}
	return &Config2{
		FromEnv:     *fromEnv,
		FromCfgFile: *fromFile,
	}, nil
}

type FromEnv struct {
	PostgresConnection string // `toml:"postgres_connection"`
	GlobalPasswordSalt string // `toml:"global_password_salt"`
	MongoDBUri         string // `toml:"mongodb_uri"`
	SecretSigningKey   string // `toml:"secret_signing_key"`
	SmtpHost           string // `toml:"smtp_host"`
	SmtpEmailLogin     string // `toml:"smtp_email_login"`
	SmtpEmailPasswd    string // `toml:"smtp_email_passwd"`
}

type FromCfgFile struct {
	ApplicationPort      *string  `toml:"application_port"`
	QueryComplexityLimit *int     `toml:"query_complexity_limit"`
	SMTPing              *SMTPing `toml:"smtp"`
	Logging              *Logging `toml:"log"`
}

type SMTPing struct {
	SMTPAuthor *string `toml:"smtp_author"`
	SMTPPort   *int    `toml:"smtp_port"`
}

type Logging struct {
	LoggingOutput *uint8  `toml:"logging_output"`
	LoggingLevel  *int8   `toml:"logging_level"`
	LoggingDBName *string `toml:"logging_db_name"`
}
