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
)

func run(w *app.Window) error {
	var ops op.Ops

	bar := window.NewBar()
	home := window.NewHome()
	origin := js.Global().Get("location").Get("origin").String()
	keltas, keltasBox := window.Loading(), window.Loading()
	go func() {
		defer w.Invalidate()

		ferries, err := schedule.Fetch(origin + "/tvarkarastis/")
		if err != nil {
			keltas, keltasBox = window.Error(err), window.Error(err)
			return
		}

		if len(ferries) != 1 {
			err = fmt.Errorf("no ferries found")
			keltas, keltasBox = window.Error(err), window.Error(err)
			return
		}

		keltas, keltasBox = window.Schedule(ferries[0]), window.Schedule(ferries[0].First(5))
	}()

	kautra, kautraBox := window.Loading(), window.Loading()
	go func() {
		defer w.Invalidate()

		bus, err := schedule.FetchBus(origin + "/lt/tvarkarastis.php")
		if err != nil {
			kautra, kautraBox = window.Error(err), window.Error(err)
			return
		}

		kautra, kautraBox = window.Bus(bus, home), window.Bus(bus.First(5), home)
	}()

	oras, orasBox := window.Loading(), window.Loading()
	go func() {
		defer w.Invalidate()

		forecasts, err := weather.Fetch(origin, weather.Places...)
		if err != nil {
			oras, orasBox = window.Error(err), window.Error(err)
			return
		}

		oras, orasBox = window.NewWeather(forecasts, home).Layout, window.NewWeather(weather.First(forecasts, 5), home).Layout
	}()

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, ui.ColorBackground)
			homePage := func(gtx layout.Context) { home.Layout(gtx, keltasBox, kautraBox, orasBox) }
			pages := []func(layout.Context){homePage, keltas, kautra, oras}
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = gtx.Constraints.Max
					pages[bar.Selected](gtx)
					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
				layout.Rigid(bar.Layout),
			)
			e.Frame(gtx.Ops)
		}
	}
}
