package main

import (
	"os"

	"github.com/vishnunotnice/NetAssist/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
