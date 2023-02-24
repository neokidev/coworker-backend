package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

type DatabaseTestSuite struct {
	suite.Suite
}

func TestDatabaseTestSuite(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		HostConfigModifier: func(config *container.HostConfig) {
			config.AutoRemove = true
		},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	testContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal("cannot create container for testing:", err)
	}

	port, err := testContainer.MappedPort(ctx, "5432")

	sourceName := fmt.Sprintf("postgresql://postgres:secret@127.0.0.1:%d/postgres?sslmode=disable", port.Int())
	testDB, err = sql.Open("postgres", sourceName)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	driver, err := postgres.WithInstance(testDB, &postgres.Config{})
	if err != nil {
		log.Fatal("cannot create migration instance:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../migration",
		"postgres", driver)
	if err != nil {
		log.Fatal("cannot create migration instance:", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatal("cannot migrate up:", err)
	}

	suite.Run(t, new(DatabaseTestSuite))
}
