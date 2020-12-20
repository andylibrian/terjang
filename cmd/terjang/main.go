package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"
	cli "github.com/urfave/cli/v2"
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
		Name:  "Terjang",
		Usage: "A scalable HTTP load testing tool built on Vegeta.",
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

					srv := server.NewServer()

					fmt.Printf("Server is listening on %s:%s\n", host, port)
					err := srv.Run(host + ":" + port)

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

					w := worker.NewWorker()
					w.SetName(name)

					fmt.Printf("Connecting to server %s:%s\n", host, port)
					w.Run(host + ":" + port)

					return nil
				},
			},
		},
	}
}
