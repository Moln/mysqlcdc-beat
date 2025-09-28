package main

import (
	"os"

	"github.com/moln/cdcbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
