package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"keltas/schedule"
	"keltas/window"
	"syscall/js"
)

func run(w *app.Window) error {
	var ops op.Ops

	bar := window.NewBar()
	bar.Selected = 1
	empty := func(gtx layout.Context) {}
	paint := window.Loading()
	go func() {
		defer w.Invalidate()

		ferries, err := schedule.Fetch(js.Global().Get("location").Get("origin").String() + "/tvarkarastis/")
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

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			pages := []func(layout.Context){empty, paint, empty, empty}
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
