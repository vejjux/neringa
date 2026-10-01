package window

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	ui "keltas/window/schedule"
)

const (
	barInset       = unit.Dp(5)
	barRadius      = unit.Dp(10)
	barButtonWidth = unit.Dp(110)
)

var barTitles = [...]string{"Home", "Naujoji", "Senoji", "Autobusai", "Orai"}

type Bar struct {
	Selected int
	buttons  [len(barTitles)]widget.Clickable
	th       *material.Theme
}

func NewBar() *Bar {
	return &Bar{th: newTheme()}
}

func (b *Bar) Layout(gtx layout.Context) layout.Dimensions {
	height := 0
	children := make([]layout.FlexChild, len(b.buttons)+1)
	for i := range b.buttons {
		if b.buttons[i].Clicked(gtx) {
			b.Selected = i
		}
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(barButtonWidth)
			gtx.Constraints.Max.X = gtx.Constraints.Min.X
			btn := material.Button(b.th, &b.buttons[i], barTitles[i])
			btn.CornerRadius = 0
			btn.Background, btn.Color = ui.ColorHeader, ui.ColorOnHeader
			if i == b.Selected {
				btn.Background, btn.Color = ui.ColorAccent, ui.ColorOnAccent
			}
			dims := btn.Layout(gtx)
			height = max(height, dims.Size.Y)
			return dims
		})
	}
	children[len(b.buttons)] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		size := image.Pt(gtx.Constraints.Max.X, height)
		paint.FillShape(gtx.Ops, ui.ColorHeader, clip.Rect{Max: size}.Op())
		return layout.Dimensions{Size: size}
	})
	return layout.UniformInset(barInset).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		m := op.Record(gtx.Ops)
		dims := layout.Flex{}.Layout(gtx, children...)
		call := m.Stop()
		defer clip.UniformRRect(image.Rectangle{Max: dims.Size}, gtx.Dp(barRadius)).Push(gtx.Ops).Pop()
		call.Add(gtx.Ops)
		return dims
	})
}
