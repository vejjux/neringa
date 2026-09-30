package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
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
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			paint(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
