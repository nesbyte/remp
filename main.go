package main

import (
	"os"
	"remp/src"
)

func main() {
	app := src.SetupCli()
	_ = app.Run(os.Args)
}
