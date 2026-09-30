package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"keltas/schedule"
	"keltas/weather"
	"keltas/window"
	"syscall/js"
)

func run(w *app.Window) error {
	var ops op.Ops

	bar := window.NewBar()
	home := window.NewHome()
	origin := js.Global().Get("location").Get("origin").String()
	paint := window.Loading()
	go func() {
		defer w.Invalidate()

		ferries, err := schedule.Fetch(origin + "/tvarkarastis/")
		if err != nil {
			paint = window.Error(err)
			return
		}

		if len(ferries) != 1 {
			paint = window.Error(fmt.Errorf("no ferries found"))
			return
		}

		paint = window.Schedule(ferries[0])
	}()

	kautra := window.Loading()
	go func() {
		defer w.Invalidate()

		bus, err := schedule.FetchBus(origin + "/lt/tvarkarastis.php")
		if err != nil {
			kautra = window.Error(err)
			return
		}

		kautra = window.Bus(bus, home)
	}()

	oras := window.Loading()
	go func() {
		defer w.Invalidate()

		forecasts, err := weather.Fetch(origin, weather.Places...)
		if err != nil {
			oras = window.Error(err)
			return
		}

		oras = window.NewWeather(forecasts, home).Layout
	}()

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			pages := []func(layout.Context){home.Layout, paint, kautra, oras}
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
