package main

import (
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	go func() {
		if err := run(app.NewWindow()); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
