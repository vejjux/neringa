package window

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	ui "keltas/window/schedule"
)

func Error(err error) func(gtx layout.Context) {
	return func(gtx layout.Context) {
		th := newTheme()
		title := material.H3(th, "Error:")
		title.Color = ui.ColorError
		title.Alignment = text.Middle
		title.Layout(gtx)

		body := material.Body1(th, err.Error())
		body.Alignment = text.Start
		body.Layout(gtx)
	}
}
