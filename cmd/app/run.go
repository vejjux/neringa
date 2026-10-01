package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"keltas/schedule"
	"keltas/weather"
	"keltas/window"
	ui "keltas/window/schedule"
	"syscall/js"
	"time"
)

func run(w *app.Window) error {
	var ops op.Ops

	bar := window.NewBar()
	home := window.NewHome()
	origin := js.Global().Get("location").Get("origin").String()
	home.Open = func(page int) { bar.Selected = page }
	loading := window.Loading()
	naujoji, naujojiTile := loading, loading
	senojiFerry, senojiTile, kautraFromSmiltyne := loading, loading, loading
	go func() {
		defer w.Invalidate()

		ferries, err := schedule.Fetch(origin + "/tvarkarastis/")
		if err == nil && len(ferries) == 0 {
			err = fmt.Errorf("no ferries found")
		}
		if err != nil {
			naujoji, senojiFerry = window.Error(err), window.Error(err)
			naujojiTile, senojiTile = window.TileError(err), window.TileError(err)
			return
		}

		if f, err := schedule.Find(ferries, "NAUJOJI"); err != nil {
			naujoji, naujojiTile = window.Error(err), window.TileError(err)
		} else {
			naujoji, naujojiTile = window.Schedule(f), window.FerryTile(f)
		}

		if f, err := schedule.Find(ferries, "SENOJI"); err != nil {
			senojiFerry, senojiTile = window.Error(err), window.TileError(err)
		} else {
			senojiFerry, senojiTile = window.Schedule(f), window.FerryTile(f)
		}
	}()

	kautra, kautraTile := loading, loading
	go func() {
		defer w.Invalidate()

		bus, err := schedule.FetchBus(origin + "/lt/tvarkarastis.php")
		if err != nil {
			kautra, kautraFromSmiltyne, kautraTile = window.Error(err), window.Error(err), window.TileError(err)
			return
		}

		kautra, kautraFromSmiltyne, kautraTile = window.Bus(bus, home), window.BusFromSmiltyne(bus), window.BusTile(bus, home)
	}()

	oras, orasTile := loading, loading
	go func() {
		defer w.Invalidate()

		forecasts, err := weather.Fetch(origin, weather.Places...)
		if err != nil {
			oras, orasTile = window.Error(err), window.TileError(err)
			return
		}

		oras, orasTile = window.NewWeather(forecasts, home).Layout, window.WeatherTile(forecasts, home)
	}()

	// Redraw periodically so the next departures on the tiles stay current.
	go func() {
		for range time.Tick(30 * time.Second) {
			w.Invalidate()
		}
	}()

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, ui.ColorBackground)
			homePage := func(gtx layout.Context) {
				home.Layout(gtx,
					window.Tile{Page: 1, Content: naujojiTile},
					window.Tile{Page: 2, Content: senojiTile},
					window.Tile{Page: 3, Content: kautraTile},
					window.Tile{Page: 4, Content: orasTile},
				)
			}
			senoji := window.Stack([]float32{2, 1}, senojiFerry, kautraFromSmiltyne)
			pages := []func(layout.Context){homePage, naujoji, senoji, kautra, oras}
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(bar.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = gtx.Constraints.Max
					pages[bar.Selected](gtx)
					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
