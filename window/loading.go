package window

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	ui "keltas/window/schedule"
)

func Loading() func(gtx layout.Context) {
	return func(gtx layout.Context) {
		th := newTheme()
		title := material.Label(th, 16, "Kraunami duomenys")
		title.Color = ui.ColorTableText
		title.Alignment = text.Middle
		title.Layout(gtx)
	}
}
