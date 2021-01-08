package main

import (
	"os"

	"github.com/yukshimizu/pulsarbeat/cmd"

	_ "github.com/yukshimizu/pulsarbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
