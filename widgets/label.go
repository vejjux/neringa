package widgets

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
)

type Label struct {
	Font       font.Font
	Color      color.NRGBA
	Background color.NRGBA
	Alignment  text.Alignment
	Text       string
	TextSize   unit.Sp
	Padding    unit.Dp
	shaper     *text.Shaper
}

func (l Label) Layout(gtx layout.Context) (dims layout.Dimensions) {
	paint.FillShape(gtx.Ops, l.Background, clip.Rect{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Sp(l.TextSize+4) + 2*gtx.Dp(l.Padding)},
	}.Op())

	m := op.Record(gtx.Ops)
	paint.ColorOp{Color: l.Color}.Add(gtx.Ops)
	textColor := m.Stop()
	tl := widget.Label{Alignment: l.Alignment, MaxLines: 1}
	dims = layout.UniformInset(l.Padding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return tl.Layout(gtx, l.shaper, l.Font, l.TextSize, l.Text, textColor)
	})

	return
}

func NewLabel(th *material.Theme, size unit.Sp, txt string) Label {
	return Label{
		Text:       txt,
		Color:      th.Palette.Fg,
		Background: th.Palette.Bg,
		TextSize:   size,
		shaper:     th.Shaper,
	}
}
