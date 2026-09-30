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
	label.Color = ColorLight
	label.Background = ColorDark
	label.Alignment = text.Middle
	label.Font.Weight = font.Black
	return label.Layout
}
