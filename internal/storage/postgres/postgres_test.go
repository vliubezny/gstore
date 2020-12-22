//+build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

var (
	db  *sql.DB
	ctx = context.Background()
	s   storage.Storage
)

func TestMain(m *testing.M) {
	shutdown := setup()
	s = New(db)

	code := m.Run()

	shutdown()
	os.Exit(code)
}

func setup() func() {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		Env:          map[string]string{"POSTGRES_PASSWORD": "root"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})
	if err != nil {
		logrus.WithError(err).Fatalf("failed to create container")
	}

	if err := c.Start(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to start container")
	}

	host, err := c.Host(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("failed to get host")
	}

	port, err := c.MappedPort(ctx, "5432")
	if err != nil {
		logrus.WithError(err).Fatal("failed to map port")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=postgres password=root dbname=postgres sslmode=disable", host, port.Int())

	_, currFile, _, ok := runtime.Caller(0)
	if !ok {
		logrus.Fatal("failed to get current file location")
	}
	migrations := filepath.Join(currFile, "../../../../scripts/migrations/postgres/")

	db = MustSetupDB(dsn, 0, 1, migrations)

	shutdownFn := func() {
		if c != nil {
			c.Terminate(ctx)
		}
	}

	return shutdownFn
}

func TestPg_GetCategories(t *testing.T) {
	categories, err := s.GetCategories(ctx)
	require.NoError(t, err)

	assert.Equal(t, []*model.Category{
		{ID: 1, Name: "Electronics"},
		{ID: 2, Name: "Computers"},
		{ID: 3, Name: "Smart Home"},
		{ID: 4, Name: "Arts & Crafts"},
		{ID: 5, Name: "Health & Household"},
		{ID: 6, Name: "Automotive"},
		{ID: 7, Name: "Pet supplies"},
		{ID: 8, Name: "Software"},
		{ID: 9, Name: "Sports & Outdoors"},
		{ID: 10, Name: "Toys and Games"},
	}, categories)
}
