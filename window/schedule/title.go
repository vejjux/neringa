package schedule

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	"keltas/widgets"
)

func Title(th *material.Theme, value string) layout.Widget {
	label := widgets.NewLabel(th, 14, value)
	label.Color = ColorOnAccent
	label.Background = ColorAccent
	label.Alignment = text.Middle
	label.Padding = 2
	label.Font.Weight = font.Black
	return label.Layout
}
