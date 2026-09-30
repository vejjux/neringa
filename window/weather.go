package window

import (
	"fmt"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"keltas/weather"
	"keltas/widgets"
	ui "keltas/window/schedule"
	"math"
)

var compass = [...]string{"Š", "ŠR", "R", "PR", "P", "PV", "V", "ŠV"}

type Weather struct {
	forecasts map[string][]weather.Hour
	home      *Home
	list      widget.List
	th        *material.Theme
}

func NewWeather(forecasts map[string][]weather.Hour, home *Home) *Weather {
	w := &Weather{forecasts: forecasts, home: home, th: newTheme()}
	w.list.Axis = layout.Vertical
	return w
}

func (w *Weather) Layout(gtx layout.Context) {
	place := StopPlaces[Stops[w.home.Selected]]
	left, right := w.forecasts[place.Code], w.forecasts[weather.Smiltyne.Code]

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(ui.Title(w.th, "ORAS")),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, ui.Direction(w.th, place.Name)),
				layout.Flexed(1, ui.Direction(w.th, weather.Smiltyne.Name)),
			)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			dims := material.List(w.th, &w.list).Layout(gtx, max(len(left), len(right)), func(gtx layout.Context, i int) layout.Dimensions {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(1, w.hour(left, i)),
					layout.Flexed(1, w.hour(right, i)),
				)
			})
			half := dims.Size.X / 2
			border(gtx, image.Pt(half, dims.Size.Y))
			defer op.Offset(image.Pt(half, 0)).Push(gtx.Ops).Pop()
			border(gtx, image.Pt(dims.Size.X-half, dims.Size.Y))
			return dims
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			credits := material.Caption(w.th, "Duomenų šaltinis: Lietuvos hidrometeorologijos tarnyba, CC BY-SA 4.0")
			credits.Alignment = text.Middle
			return layout.UniformInset(4).Layout(gtx, credits.Layout)
		}),
	)
}

func border(gtx layout.Context, size image.Point) {
	gtx.Constraints = layout.Exact(size)
	widget.Border{Color: color.NRGBA{A: 0xFF}, Width: 1}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: size}
	})
}

func (w *Weather) hour(hours []weather.Hour, i int) layout.Widget {
	if i >= len(hours) {
		return func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{} }
	}
	h := hours[i]

	hour := widgets.NewLabel(w.th, 20, fmt.Sprintf(" %s: ", h.Time.Local().Format("15")))
	hour.Color = ui.ColorDark
	hour.Font.Weight = font.Bold
	hour.Alignment = text.Start

	value := fmt.Sprintf("%d° %s %d m/s", int(math.Round(h.Temperature)), compass[int(math.Round(h.WindDirection/45))%8], int(math.Round(h.WindSpeed)))
	if h.Precipitation > 0 {
		value += fmt.Sprintf(" %.1f mm", h.Precipitation)
	}
	values := widgets.NewLabel(w.th, 16, value)
	values.Color = ui.ColorDark
	values.Font.Weight = font.Bold
	values.Alignment = text.Start

	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx, layout.Rigid(hour.Layout), layout.Rigid(values.Layout))
	}
}
