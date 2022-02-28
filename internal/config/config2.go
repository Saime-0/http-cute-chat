package config

import (
	"github.com/BurntSushi/toml"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
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
		return nil, cerrors.Wrap(err, "не удалось декодировать файл")
	}

	if !fromFile.validate() {
		return nil, cerrors.New("в файле конфигурации заполнены не все поля")
	}
	if !clog.Exists(clog.LogLevel(*fromFile.Logging.LoggingLevel)) {
		return nil, cerrors.New("указан несуществующий уровень лога")
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
		return nil, cerrors.New("не установлены некоторые переменные окружения")
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
	ApplicationPort                   *string  `toml:"application_port"`
	QueryComplexityLimit              *int     `toml:"query_complexity_limit"`
	DurationOfScheduleInterval        *int64   `toml:"duration_of_schedule_interval"`
	RefreshTokenLiftime               *int64   `toml:"refresh_token_liftime"`
	AccessTokenLiftime                *int64   `toml:"access_token_liftime"`
	MaximumNumberOfMessagesPerRequest *int     `toml:"maximum_number_of_messages_per_request"`
	MaxCountRooms                     *int     `toml:"max_count_rooms"`
	MaxUserChats                      *int     `toml:"max_user_chats"`
	MaxCountOwnedChats                *int     `toml:"max_count_owned_chats"`
	MaxMembersOnChat                  *int     `toml:"max_members_on_chat"`
	MaxRolesInChat                    *int     `toml:"max_roles_in_chat"`
	SMTPing                           *SMTPing `toml:"smtp"`
	Logging                           *Logging `toml:"log"`
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
