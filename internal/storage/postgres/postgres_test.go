//+build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vliubezny/gstore/internal/storage"
)

type postgresTestSuite struct {
	suite.Suite
	ctx       context.Context
	container testcontainers.Container
	db        *sql.DB
	migrator  *migrate.Migrate
	s         storage.Storage
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, &postgresTestSuite{})
}

func (s *postgresTestSuite) SetupSuite() {
	s.ctx = context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		Env:          map[string]string{"POSTGRES_PASSWORD": "root"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	c, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})
	s.Require().NoError(err, "failed to create container")
	s.container = c

	s.Require().NoError(c.Start(s.ctx), "failed to start container")

	host, err := c.Host(s.ctx)
	s.Require().NoError(err, "failed to get host")

	port, err := c.MappedPort(s.ctx, "5432")
	s.Require().NoError(err, "failed to map port")

	dsn := fmt.Sprintf("host=%s port=%d user=postgres password=root dbname=postgres sslmode=disable", host, port.Int())
	s.T().Log(dsn)
	s.db, s.migrator = MustPrepareDB(dsn, 0, 1, "../../../scripts/migrations/postgres/")

	s.s = New(s.db)
}

func (s *postgresTestSuite) TearDownSuite() {
	if s.container != nil {
		s.container.Terminate(s.ctx)
	}
}

func (s *postgresTestSuite) SetupTest() {
	if err := s.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}
}

func (s *postgresTestSuite) TearDownTest() {
	s.NoError(s.migrator.Down())
}
