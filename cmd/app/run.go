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
	naujoji, naujojiBox := window.Loading(), window.Loading()
	senojiFerry, kautraFromSmiltyne := window.Loading(), window.Loading()
	go func() {
		defer w.Invalidate()

		ferries, err := schedule.Fetch(origin + "/tvarkarastis/")
		if err == nil && len(ferries) == 0 {
			err = fmt.Errorf("no ferries found")
		}
		if err != nil {
			naujoji, naujojiBox, senojiFerry = window.Error(err), window.Error(err), window.Error(err)
			return
		}

		if f, err := schedule.Find(ferries, "NAUJOJI"); err != nil {
			naujoji, naujojiBox = window.Error(err), window.Error(err)
		} else {
			naujoji, naujojiBox = window.Schedule(f), func(gtx layout.Context) {
				window.Schedule(f.Upcoming(time.Now().Hour(), 5))(gtx)
			}
		}

		if f, err := schedule.Find(ferries, "SENOJI"); err != nil {
			senojiFerry = window.Error(err)
		} else {
			senojiFerry = window.Schedule(f)
		}
	}()

	kautra, kautraBox := window.Loading(), window.Loading()
	go func() {
		defer w.Invalidate()

		bus, err := schedule.FetchBus(origin + "/lt/tvarkarastis.php")
		if err != nil {
			kautra, kautraBox, kautraFromSmiltyne = window.Error(err), window.Error(err), window.Error(err)
			return
		}

		kautra, kautraBox = window.Bus(bus, home), func(gtx layout.Context) {
			window.Bus(bus.Upcoming(time.Now().Hour(), 5), home)(gtx)
		}
		kautraFromSmiltyne = window.BusFromSmiltyne(bus)
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
			homePage := func(gtx layout.Context) { home.Layout(gtx, naujojiBox, kautraBox, orasBox) }
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
