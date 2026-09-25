package main

import (
	"github.com/umatare5/twelvedata-exporter/cli"
	"github.com/umatare5/twelvedata-exporter/log"
)

// The entrypoint of this program.
func main() {
	if err := cli.Start(); err != nil {
		log.Fatal(err)
	}
}
