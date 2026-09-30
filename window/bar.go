package window

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	ui "keltas/window/schedule"
)

var barTitles = [...]string{"Home", "Keltas", "Kautra", "Oras"}

type Bar struct {
	Selected int
	buttons  [len(barTitles)]widget.Clickable
	th       *material.Theme
}

func NewBar() *Bar {
	return &Bar{th: newTheme()}
}

func (b *Bar) Layout(gtx layout.Context) layout.Dimensions {
	children := make([]layout.FlexChild, len(b.buttons))
	for i := range b.buttons {
		if b.buttons[i].Clicked(gtx) {
			b.Selected = i
		}
		children[i] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(b.th, &b.buttons[i], barTitles[i])
			btn.CornerRadius = 0
			btn.Background, btn.Color = ui.ColorHeader, ui.ColorOnHeader
			if i == b.Selected {
				btn.Background, btn.Color = ui.ColorAccent, ui.ColorOnAccent
			}
			return btn.Layout(gtx)
		})
	}
	return layout.Flex{}.Layout(gtx, children...)
}
