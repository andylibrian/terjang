package main

import (
	"fmt"
	"log"
	"os"

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

					fmt.Printf("Server is listening on %s:%s\n", host, port)

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
					host := c.String("host")
					port := c.String("port")

					fmt.Printf("Connecting to server %s:%s\n", host, port)

					return nil
				},
			},
		},
	}
}
