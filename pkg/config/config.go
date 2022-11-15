package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	// HttpPort - http server port
	HttpPort = "PORT"
	HttpHost = "HOST"

	// DbHost - db host key
	DbHost = "DB_HOST"
	// DbPort - db port key
	DbPort = "DB_PORT"
	// DbUser - db user key
	DbUser = "DB_USER"
	// DbPass - db pass key
	DbPass = "DB_PASS"
	// DbName - db name key
	DbName = "DB_NAME"
	// DbSsl - db ssl mode key [disable/enable]
	DbSsl = "DB_SSL"

	// LogLevel - logging level
	LogLevel = "LOG_LEVEL"
)

type Config struct {
	DbCredentials *DatabaseCredentials
	Http          *Http
	Logger        *Logger
}

type Logger struct {
	Level logrus.Level
}

type Http struct {
	Port string
	Host string
}

type DatabaseCredentials struct {
	Host        string
	Port        string
	Name        string
	User        string
	Pass        string
	SslMode     string
	AutoMigrate bool
}

type RedisCredentials struct {
	Host string
	Port string
	Pass string
}

func InitDbCredentials() *DatabaseCredentials {
	return &DatabaseCredentials{
		Host:        os.Getenv(DbHost),
		Name:        os.Getenv(DbName),
		Port:        os.Getenv(DbPort),
		Pass:        os.Getenv(DbPass),
		User:        os.Getenv(DbUser),
		SslMode:     os.Getenv(DbSsl),
		AutoMigrate: false,
	}
}

// InitHttpConfig - return pointer of http server config
func InitHttpConfig() *Http {
	return &Http{
		Port: os.Getenv(HttpPort),
		Host: os.Getenv(HttpHost),
	}
}

// InitLoggerConfig - return pointer of Logger config
func InitLoggerConfig() *Logger {
	level := os.Getenv(LogLevel)

	var lvl logrus.Level

	switch level {
	case "debug":
		lvl = logrus.DebugLevel
	case "info":
		lvl = logrus.InfoLevel
	case "warn":
		lvl = logrus.WarnLevel
	case "error":
		lvl = logrus.ErrorLevel
	case "fatal":
		lvl = logrus.FatalLevel
	default:
		lvl = logrus.FatalLevel
	}

	return &Logger{
		Level: lvl,
	}
}

func Init() *Config {
	return &Config{
		DbCredentials: InitDbCredentials(),
		Http:          InitHttpConfig(),
		Logger:        InitLoggerConfig(),
	}
}
