package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"github.com/vliubezny/gstore/internal/auth"
	"github.com/vliubezny/gstore/internal/server"
	"github.com/vliubezny/gstore/internal/service"
	"github.com/vliubezny/gstore/internal/storage"
	"github.com/vliubezny/gstore/internal/storage/postgres"
	"golang.org/x/sync/errgroup"
)

var errTerminated = errors.New("terminated")

var opts = struct {
	Host string `long:"http.host" env:"HTTP_HOST" default:"0.0.0.0" description:"IP address to listen"`
	Port int    `long:"http.port" env:"HTTP_PORT" default:"8080" description:"port to listen"`

	LogLevel string `long:"log.level" env:"LOG_LEVEL" default:"debug" description:"Log level" choice:"debug" choice:"info" choice:"warning" choice:"error"`

	SignKey string `long:"auth.signkey" env:"AUTH_SIGN_KEY" default:"changeme" description:"sign key for JWT"`

	PostgresDSN                string `long:"postgres" env:"POSTGRES_DSN" default:"host=localhost port=5432 user=postgres password=root dbname=postgres sslmode=disable" description:"postgres dsn"`
	PostgresMaxOpenConnections int    `long:"postgres.max_open_connections" env:"POSTGRES_MAX_OPEN_CONNECTIONS" default:"0" description:"postgres maximal open connections count, 0 means unlimited"`
	PostgresMaxIdleConnections int    `long:"postgres.max_idle_connections" env:"POSTGRES_MAX_IDLE_CONNECTIONS" default:"5" description:"postgres maximal idle connections count"`
	PostgresMigrations         string `long:"postgres.migrations" env:"POSTGRES_MIGRATIONS" default:"scripts/migrations/postgres" description:"postgres migrations directory"`
}{}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gstore"
	parser.LongDescription = "Starts gstore server."

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		logrus.WithError(err).Fatal("failed to parse flags")
	}

	lvl, _ := logrus.ParseLevel(opts.LogLevel)
	logrus.SetLevel(lvl)

	logrus.Info("starting service")
	logrus.Infof("%+v", opts) // can print secrets!

	db := postgres.MustSetupDB(opts.PostgresDSN, opts.PostgresMaxOpenConnections,
		opts.PostgresMaxIdleConnections, opts.PostgresMigrations)
	strg := postgres.New(db)
	authSvc := auth.New(strg.(storage.UserStorage), opts.SignKey)
	r := chi.NewMux()

	server.SetupRouter(service.New(strg), authSvc, r, authSvc.ValidateAccessToken)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler: r,
	}

	gr, _ := errgroup.WithContext(context.Background())
	gr.Go(srv.ListenAndServe)

	gr.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		s := <-sigs
		logrus.Infof("terminating by %s signal", s)

		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.WithError(err).Error("failed to gracefully shutdown server")
		}

		return errTerminated
	})

	logrus.Info("service started")

	if err := gr.Wait(); err != nil && !errors.Is(err, errTerminated) && !errors.Is(err, http.ErrServerClosed) {
		logrus.WithError(err).Fatal("service unexpectedly stopped")
	}
}
