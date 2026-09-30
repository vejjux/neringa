package main

import (
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Neringa"))
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
