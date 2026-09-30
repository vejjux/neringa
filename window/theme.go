package window

import (
	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget/material"
	ui "keltas/window/schedule"
)

var theme *material.Theme

func newTheme() *material.Theme {
	if theme == nil {
		theme = material.NewTheme()
		theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
		theme.Palette = material.Palette{
			Bg:         ui.ColorBackground,
			Fg:         ui.ColorText,
			ContrastBg: ui.ColorAccent,
			ContrastFg: ui.ColorOnAccent,
		}
	}
	return theme
}
