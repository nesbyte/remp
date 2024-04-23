package main

import (
	"os"
	"remp/src"
)

var Version = "development"

func main() {
	app := src.SetupCli(Version)
	_ = app.Run(os.Args)
}
