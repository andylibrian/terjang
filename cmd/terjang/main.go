package main

import (
	"log"
	"os"

	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = "main"
)

func main() {
	app := getCliApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getCliApp() *cli.App {
	return &cli.App{
		Name:    "Terjang",
		Usage:   "A scalable HTTP load testing tool built on Vegeta.",
		Version: version + " (" + commit + ")",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "Log level: debug, info, warn, error.",
				Value: "info",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Run server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Usage: "Host address to listen on",
						Value: "0.0.0.0",
					},
					&cli.StringFlag{
						Name:  "port",
						Usage: "Host port to listen on",
						Value: "9009",
					},
				},
				Action: func(c *cli.Context) error {
					host := c.String("host")
					port := c.String("port")
					logLevel := c.String("log-level")

					logger := getLogger(logLevel)
					server.SetLogger(logger)

					srv := server.NewServer()

					err := srv.Run(host + ":" + port)
					defer srv.Close()

					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "worker",
				Usage: "Run worker",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Usage:       "name of the worker",
						DefaultText: "hostname",
					},
					&cli.StringFlag{
						Name:  "host",
						Usage: "Sever's host address to connect to",
						Value: "localhost",
					},
					&cli.StringFlag{
						Name:  "port",
						Usage: "Server's host port to connect to",
						Value: "9009",
					},
				},
				Action: func(c *cli.Context) error {
					name := c.String("name")
					if name == "" {
						hostname, err := os.Hostname()

						if err == nil {
							name = hostname
						} else {
							name = "worker"
						}
					}
					host := c.String("host")
					port := c.String("port")
					logLevel := c.String("log-level")

					logger := getLogger(logLevel)
					worker.SetLogger(logger)

					w := worker.NewWorker()
					w.SetName(name)

					w.Run(host + ":" + port)

					return nil
				},
			},
		},
	}
}

func getLogger(level string) *zap.SugaredLogger {
	zapLevel := zap.InfoLevel
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()

	if err != nil {
		log.Fatal(err)
	}

	return logger.Sugar()
}
