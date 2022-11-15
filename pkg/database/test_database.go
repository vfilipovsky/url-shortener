package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	DB  *sql.DB
	URI string
}

func SetupTestDatabase() (*TestDatabase, error) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:13.5-alpine",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		},
	)

	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	URI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("postgres", URI)

	if err != nil {
		return nil, err
	}

	return &TestDatabase{DB: db, URI: URI}, nil
}

func (t *TestDatabase) Migrate(up bool) error {
	var err error

	driver, err := postgres.WithInstance(t.DB, &postgres.Config{})

	if err != nil {
		println(err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/migrations", rootDir()), "postgres", driver)

	if err != nil {
		println(err)
		os.Exit(1)
	}

	if up {
		err = m.Up()
	} else {
		err = m.Down()
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func rootDir() string {
	f, err := os.Getwd()

	if err != nil {
		println(err)
		os.Exit(1)
	}

	splitted := strings.Split(filepath.Dir(f), "/")

	for len(splitted) > 0 && splitted[len(splitted)-1] != "url-shortener" {
		splitted = splitted[:len(splitted)-1]
	}

	return strings.Join(splitted, "/")
}
