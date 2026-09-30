package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/op"
	"keltas/schedule"
	"keltas/window"
	"syscall/js"
)

func run(w *app.Window) error {
	var ops op.Ops

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
			paint(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
