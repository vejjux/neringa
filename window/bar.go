package window

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"

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
	barHeight      = unit.Dp(42)
	barButtonWidth = unit.Dp(110)
	barIconWidth   = unit.Dp(56)
	barIconInset   = unit.Dp(6)
	barPadding     = unit.Dp(2)
)

// barTitles names the pages; the first one, home, is shown as the app icon.
var barTitles = [...]string{"Home", "Naujoji", "Senoji", "Autobusai", "Orai"}

//go:embed icon.png
var iconPNG []byte

var icon = func() paint.ImageOp {
	img, err := png.Decode(bytes.NewReader(iconPNG))
	if err != nil {
		panic(err)
	}
	return paint.NewImageOp(img)
}()

type Bar struct {
	Selected int
	buttons  [len(barTitles)]widget.Clickable
	th       *material.Theme
}

func NewBar() *Bar {
	return &Bar{th: newTheme()}
}

func (b *Bar) Layout(gtx layout.Context) layout.Dimensions {
	for i := range b.buttons {
		if b.buttons[i].Clicked(gtx) {
			b.Selected = i
		}
	}
	return layout.UniformInset(barInset).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		widths := b.widths(gtx)
		children := make([]layout.FlexChild, len(b.buttons)+1)
		for i := range b.buttons {
			children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints = layout.Exact(image.Pt(widths[i], gtx.Dp(barHeight)))
				background, color := ui.ColorHeader, ui.ColorOnHeader
				if i == b.Selected {
					background, color = ui.ColorAccent, ui.ColorOnAccent
				}
				if i == 0 {
					btn := material.ButtonLayout(b.th, &b.buttons[i])
					btn.CornerRadius = 0
					btn.Background = background
					return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(barIconInset).Layout(gtx, widget.Image{Src: icon, Fit: widget.Contain}.Layout)
					})
				}
				btn := material.Button(b.th, &b.buttons[i], barTitles[i])
				btn.CornerRadius = 0
				btn.Inset.Left, btn.Inset.Right = barPadding, barPadding
				btn.Background, btn.Color = background, color
				return btn.Layout(gtx)
			})
		}
		children[len(b.buttons)] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Dp(barHeight))}
		})

		m := op.Record(gtx.Ops)
		dims := layout.Flex{}.Layout(gtx, children...)
		call := m.Stop()
		defer clip.UniformRRect(image.Rectangle{Max: dims.Size}, gtx.Dp(barRadius)).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, ui.ColorHeader)
		call.Add(gtx.Ops)
		return dims
	})
}

// widths gives the icon its fixed width and shares the rest of the bar
// equally between the text buttons, up to barButtonWidth each.
func (b *Bar) widths(gtx layout.Context) []int {
	widths := make([]int, len(barTitles))
	widths[0] = gtx.Dp(barIconWidth)
	w := min(gtx.Dp(barButtonWidth), (gtx.Constraints.Max.X-widths[0])/(len(barTitles)-1))
	for i := 1; i < len(widths); i++ {
		widths[i] = w
	}
	return widths
}
