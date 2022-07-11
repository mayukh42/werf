package main

import (
	"flag"
	"log"

	"github.com/mayukh42/werf/app/werf"
	"github.com/mayukh42/werf/config"
)

const (
	CONFIG_FILE_LOCATION = "./resources/dev"
)

func main() {
	// run
	cfgFile := flag.String("config", CONFIG_FILE_LOCATION, "config file location")
	log.Printf("using config file from %s", *cfgFile)

	quayName := flag.String("name", "main", "default name: main")
	id := flag.Int("id", 0, "default id: 0")
	log.Printf("using quay name %s and id %d", *quayName, *id)

	cfg := config.GetConfig(&config.ConfigInput{
		Name: "config",
		Type: "yml",
		Path: *cfgFile,
	})

	q := werf.NewQuay(cfg)
	q.Run(*quayName, *id)
}
