package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	filePath := "../../.env.test"
	err := godotenv.Load(filePath)

	if err != nil {
		panic(fmt.Sprintf("'%s' file not found", filePath))
	}

	os.Exit(m.Run())
}

func TestInitLoggerConfig(t *testing.T) {
	expected := &Logger{
		Level: logrus.ErrorLevel,
	}

	actual := InitLoggerConfig()

	assert.EqualValues(t, *expected, *actual)
}

func TestInitDatabaseCredentials(t *testing.T) {
	expected := &DatabaseCredentials{
		Name:        "testdb",
		Host:        "127.0.0.1",
		Port:        "5432",
		User:        "postgres",
		Pass:        "postgres",
		SslMode:     "disable",
		AutoMigrate: false,
	}

	actual := InitDbCredentials()

	assert.Equal(t, *expected, *actual)
}

func TestInitHttpConfig(t *testing.T) {
	expected := &Http{
		Port: "8085",
		Host: "127.0.0.1",
	}

	actual := InitHttpConfig()

	assert.Equal(t, *expected, *actual)
}

func TestInit(t *testing.T) {
	dbCredentials := &DatabaseCredentials{
		Name:        "testdb",
		Host:        "127.0.0.1",
		Port:        "5432",
		User:        "postgres",
		Pass:        "postgres",
		SslMode:     "disable",
		AutoMigrate: false,
	}

	httpConfig := &Http{
		Port: "8085",
		Host: "127.0.0.1",
	}

	loggerConfig := &Logger{
		Level: logrus.ErrorLevel,
	}

	expected := &Config{
		DbCredentials: dbCredentials,
		Http:          httpConfig,
		Logger:        loggerConfig,
	}

	actual := Init()

	assert.Equal(t, *expected, *actual)
}
