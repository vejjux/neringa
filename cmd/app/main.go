package main

import (
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	go func() {
		if err := run(new(app.Window)); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
