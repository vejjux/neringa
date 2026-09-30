package window

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"keltas/storage"
	ui "keltas/window/schedule"
)

const homeKey = "mano-namai"

var Stops = [...]string{
	"Nidos gyvenvietės autobusų stotis",
	"G. D. Kuverto plento sankryža",
	"T. Mano muziejus",
	"Preilos gv.",
	"Preilos gv. prie plento",
	"Pervalkos gv.",
	"Pervalkos gv. prie plento",
	"Žvejų kaimelis",
	"Raganų kalnas",
	"Juodkrantės gv.",
	"Gintaro įlanka",
	"Alksnynė",
}

type Home struct {
	Selected int
	open     bool
	toggle   widget.Clickable
	options  [len(Stops)]widget.Clickable
	list     widget.List
	th       *material.Theme
}

func NewHome() *Home {
	h := &Home{th: newTheme()}
	h.list.Axis = layout.Vertical
	saved := storage.Get(homeKey)
	for i, s := range Stops {
		if s == saved {
			h.Selected = i
		}
	}
	return h
}

func (h *Home) Layout(gtx layout.Context, sections ...func(layout.Context)) {
	if h.toggle.Clicked(gtx) {
		h.open = !h.open
	}
	for i := range h.options {
		if h.options[i].Clicked(gtx) {
			h.Selected, h.open = i, false
			storage.Set(homeKey, Stops[i])
		}
	}

	arrow := " ▼"
	if h.open {
		arrow = " ▲"
	}
	button := material.Button(h.th, &h.toggle, Stops[h.Selected]+arrow)
	button.CornerRadius = 0
	button.Background, button.Color = ui.ColorDark, ui.ColorLight

	children := []layout.FlexChild{layout.Rigid(button.Layout)}
	if !h.open {
		for _, section := range sections {
			children = append(children, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: 4}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = gtx.Constraints.Max
					section(gtx)
					return layout.Dimensions{Size: gtx.Constraints.Max}
				})
			}))
		}
		layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		return
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx, append(children,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.List(h.th, &h.list).Layout(gtx, len(Stops), func(gtx layout.Context, i int) layout.Dimensions {
				option := material.Button(h.th, &h.options[i], Stops[i])
				option.CornerRadius = 0
				option.Background, option.Color = ui.ColorLight, ui.ColorDark
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return option.Layout(gtx)
			})
		}),
	)...)
}
