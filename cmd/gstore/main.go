package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

var opts = struct {
	Host string `long:"http.host" env:"HTTP_HOST" default:"0.0.0.0" description:"IP address to listen"`
	Port int    `long:"http.port" env:"HTTP_PORT" default:"8080" description:"port to listen"`

	LogLevel string `long:"log.level" env:"LOG_LEVEL" default:"debug" description:"Log level" choice:"debug" choice:"info" choice:"warning" choice:"error"`
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
		logrus.WithError(err).Fatal("fail to parse flags")
	}

	lvl, _ := logrus.ParseLevel(opts.LogLevel)
	logrus.SetLevel(lvl)

	logrus.Info("Starting server")
	logrus.Infof("%+v", opts) // can print secrets!

	r := chi.NewMux()

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler: r,
	}

	srv.ListenAndServe()
}
