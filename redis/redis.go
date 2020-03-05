package redis

import (
	"github.com/philippgille/gokv/syncmap"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Create(c *cli.Context) (*Server, error) {
	options := syncmap.DefaultOptions
	store := syncmap.NewStore(options)
	return New(store, log.New(os.Stdout, "[server]", log.Flags())), nil
}
