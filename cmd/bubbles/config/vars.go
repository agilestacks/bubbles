package config

const (
	BlobMaxMemory = 1024 * 1024
)

var (
	Verbose bool
	Debug   bool
	Trace   bool

	HttpPort int

	BubblesApiSecret string
)

func Update() {
	if Trace {
		Debug = true
	}
	if Debug {
		Verbose = true
	}
}
