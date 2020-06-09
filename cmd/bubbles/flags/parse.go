package flags

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/agilestacks/bubbles/cmd/bubbles/config"
)

func Parse() {
	var apiSecretEnvVar string

	flag.BoolVar(&config.Verbose, "verbose", true, "Print progress if set")
	flag.BoolVar(&config.Debug, "debug", false, "Print debug information if set")
	flag.BoolVar(&config.Trace, "trace", false, "Print detailed trace if set")

	flag.IntVar(&config.HttpPort, "http_port", 8005, "HTTP API port to listen")
	flag.StringVar(&apiSecretEnvVar, "api_secret_env", "BUBBLES_API_SECRET",
		"Environment variable to get secret from to protect Bubbles write HTTP API, set to \"\" to disable")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr,
			`Usage:
  bubbles -http_port 80

Flags:
`)
		flag.PrintDefaults()
	}

	flag.Parse()
	config.BubblesApiSecret = lookupEnv(apiSecretEnvVar, "api_secret_env")
	config.Update()
}

func lookupEnv(envVar string, param string) string {
	if envVar == "" {
		return ""
	}
	value, exist := os.LookupEnv(envVar)
	if !exist {
		log.Fatalf("`-%s %s` is set but variable not found in process environment, try `bubbles -h`",
			param, envVar)
	}
	return value
}
