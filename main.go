package main

import (
	"github.com/tdakkota/gokv-server/redis"
	"github.com/tidwall/redcon"
	"github.com/urfave/cli/v2"
	"log"
	_ "net/http/pprof"
	"os"
)

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "bind",
			Usage:       "address to bind",
			EnvVars:     []string{"BIND"},
			Required:    false,
			DefaultText: "localhost:5000",
		},
	}
}

func runner(c *cli.Context) error {
	server, err := redis.Create(c)
	if err != nil {
		return err
	}

	bind := c.String("bind")
	if bind == "" {
		bind = "localhost:5000"
	}

	log.Println("Staring server on", bind)
	return redcon.ListenAndServe(bind, server.Handler, server.Accept, server.Closed)
}

func main() {
	app := &cli.App{
		Name:   "gokv-server",
		Usage:  "runs gokv-server",
		Action: runner,
		Flags:  flags(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal()
	}
}
