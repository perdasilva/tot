package main

import (
	"os"
	"tot/cmd/root"
)

func main() {
	cmd := root.NewCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
