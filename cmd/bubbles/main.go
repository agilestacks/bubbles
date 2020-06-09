package main

import (
	"log"

	"github.com/agilestacks/bubbles/cmd/bubbles/api"
	"github.com/agilestacks/bubbles/cmd/bubbles/config"
	"github.com/agilestacks/bubbles/cmd/bubbles/flags"
)

func main() {
	flags.Parse()
	api.Init()
	api.Listen("0.0.0.0", config.HttpPort)
	if config.Verbose {
		log.Printf("Bubbles started on HTTP port %d", config.HttpPort)
	}
	select {}
}
