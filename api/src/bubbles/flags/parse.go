package flags

import (
	"flag"
	"fmt"
	"os"

	"bubbles/config"
)

func Parse() {
	flag.BoolVar(&config.Verbose, "verbose", true, "Print progress if set")
	flag.BoolVar(&config.Debug, "debug", false, "Print debug information if set")
	flag.BoolVar(&config.Trace, "trace", false, "Print detailed trace if set")

	flag.IntVar(&config.HttpPort, "http_port", 8005, "HTTP API port to listen")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr,
			`Usage:
  bubbles -http_port 80

Flags:
`)
		flag.PrintDefaults()
	}

	flag.Parse()
	config.Update()
}
