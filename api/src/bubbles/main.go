package main

import (
	"log"

	"bubbles/api"
	"bubbles/config"
	"bubbles/flags"
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
