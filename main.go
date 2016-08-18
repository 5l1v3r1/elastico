package main

import (
	"os"

	"github.com/dutchcoders/elastico/cmd"
)

func main() {
	app := cmd.NewApp()
	app.Run(os.Args)
}
