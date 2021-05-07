package main

import (
	"context"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/sevigo/hokan/cmd/hokan/config"
)

func main() {
	conf, err := config.Environ()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatal("main: invalid configuration")
	}

	initLogger(conf.Logging)
	ctx := context.Background()

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Debug("A group of walrus emerges from the ocean")
	logrus.WithField("foo", "bar").Error("XXX")

	app, err := InitializeApplication(ctx, conf)
	if err != nil {
		logrus.Fatal("main: cannot initialize server")
	}

	logrus.WithField("host", app.server.Host).Info("InitializeApplication OK")

	g := errgroup.Group{}
	g.Go(func() error {
		logrus.WithFields(
			logrus.Fields{
				"port": conf.Server.Port,
			},
		).Info("starting the http server")
		return app.server.ListenAndServe(ctx)
	})

	if err := g.Wait(); err != nil {
		logrus.Fatal("main: program terminated")
	}
}

func initLogger(conf config.Logging) {
	var logLevel = logrus.InfoLevel
	if conf.Debug {
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetOutput(colorable.NewColorableStdout())
	// logrus.SetReportCaller(false)

	if conf.Text {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:            conf.Color,
			DisableColors:          !conf.Color,
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			DisableLevelTruncation: true,
			ForceQuote:             false,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: conf.Pretty,
		})
	}
}
