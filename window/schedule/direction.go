package schedule

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	"keltas/widgets"
)

func Direction(th *material.Theme, value string) layout.Widget {
	label := widgets.NewLabel(th, 14, value)
	label.Color = ColorDark
	label.Background = ColorLight
	label.Alignment = text.Middle
	label.Font.Weight = font.Black
	return label.Layout
}
