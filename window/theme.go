package window

import (
	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget/material"
)

var theme *material.Theme

func newTheme() *material.Theme {
	if theme == nil {
		theme = material.NewTheme()
		theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	}
	return theme
}
