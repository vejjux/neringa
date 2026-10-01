package window

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"keltas/schedule"
	"keltas/storage"
	ui "keltas/window/schedule"
)

const (
	homeKey       = "mano-namai"
	dropdownWidth = unit.Dp(180)
	homeGap       = unit.Dp(10)
	tileMinSize   = unit.Dp(170)
	tileRadius    = unit.Dp(12)
	tilePadding   = unit.Dp(12)
)

var Stops = [...]string{
	schedule.Nida,
	"G. D. Kuverto plento sankryža",
	"T. Mano muziejus",
	"Preila",
	"Preila prie plento",
	"Pervalka",
	"Pervalka prie plento",
	"Žvejų kaimelis",
	"Raganų kalnas",
	"Juodkrantė",
	"Gintaro įlanka",
	"Alksnynė",
}

type Home struct {
	Selected int
	// Open is called with a tile's page when the tile is clicked.
	Open    func(page int)
	open    bool
	toggle  widget.Clickable
	options [len(Stops)]widget.Clickable
	tiles   []widget.Clickable
	list    widget.List
	grid    widget.List
	th      *material.Theme
}

func NewHome() *Home {
	h := &Home{th: newTheme()}
	h.list.Axis = layout.Vertical
	h.grid.Axis = layout.Vertical
	saved := schedule.StopName(storage.Get(homeKey))
	for i, s := range Stops {
		if s == saved {
			h.Selected = i
		}
	}
	return h
}

func (h *Home) Layout(gtx layout.Context, tiles ...Tile) {
	if h.toggle.Clicked(gtx) {
		h.open = !h.open
	}
	for i := range h.options {
		if h.options[i].Clicked(gtx) {
			h.Selected, h.open = i, false
			storage.Set(homeKey, Stops[i])
		}
	}
	for len(h.tiles) < len(tiles) {
		h.tiles = append(h.tiles, widget.Clickable{})
	}
	for i := range tiles {
		if h.tiles[i].Clicked(gtx) && h.Open != nil {
			h.Open(tiles[i].Page)
			gtx.Execute(op.InvalidateCmd{})
		}
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: barInset, Right: barInset, Bottom: homeGap}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx, layout.Rigid(h.dropdown))
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if h.open {
				return h.optionList(gtx)
			}
			return h.layoutGrid(gtx, tiles)
		}),
	)
}

func (h *Home) dropdown(gtx layout.Context) layout.Dimensions {
	arrow := "▼"
	if h.open {
		arrow = "▲"
	}
	gtx.Constraints.Min.X = gtx.Dp(dropdownWidth)
	gtx.Constraints.Max.X = gtx.Constraints.Min.X
	btn := material.ButtonLayout(h.th, &h.toggle)
	btn.CornerRadius = barRadius
	btn.Background = ui.ColorAccent
	return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: 6, Bottom: 6, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return tileLabel(gtx, 11, font.Normal, ui.ColorOnAccent, "Gyvenvietė")
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return tileLabel(gtx, 14, font.Bold, ui.ColorOnAccent, Stops[h.Selected])
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return tileLabel(gtx, 14, font.Normal, ui.ColorOnAccent, arrow)
				}),
			)
		})
	})
}

func (h *Home) optionList(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Left: barInset}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(dropdownWidth))
		return material.List(h.th, &h.list).Layout(gtx, len(Stops), func(gtx layout.Context, i int) layout.Dimensions {
			option := material.Button(h.th, &h.options[i], Stops[i])
			option.TextSize = 14
			option.CornerRadius = 0
			option.Background, option.Color = ui.ColorHeader, ui.ColorOnHeader
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return option.Layout(gtx)
		})
	})
}

// layoutGrid lays tiles out as squares, as many per row as fit tileMinSize.
func (h *Home) layoutGrid(gtx layout.Context, tiles []Tile) layout.Dimensions {
	gap, inset := gtx.Dp(homeGap), gtx.Dp(barInset)
	width := gtx.Constraints.Max.X - 2*inset
	cols := max(1, (width+gap)/(gtx.Dp(tileMinSize)+gap))
	side := (width - (cols-1)*gap) / cols
	rows := (len(tiles) + cols - 1) / cols

	return material.List(h.th, &h.grid).Layout(gtx, rows, func(gtx layout.Context, r int) layout.Dimensions {
		for c := 0; c < cols && r*cols+c < len(tiles); c++ {
			i := r*cols + c
			stack := op.Offset(image.Pt(inset+c*(side+gap), 0)).Push(gtx.Ops)
			gtx := gtx
			gtx.Constraints = layout.Exact(image.Pt(side, side))
			h.tiles[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return h.tile(gtx, tiles[i])
			})
			stack.Pop()
		}
		return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, side+gap)}
	})
}

func (h *Home) tile(gtx layout.Context, t Tile) layout.Dimensions {
	size := gtx.Constraints.Max
	defer clip.UniformRRect(image.Rectangle{Max: size}, gtx.Dp(tileRadius)).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, ui.ColorHeader)
	layout.UniformInset(tilePadding).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: 6}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return tileLabel(gtx, 16, font.Black, ui.ColorOnAccent, barTitles[t.Page])
				})
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				t.Content(gtx)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			}),
		)
	})
	return layout.Dimensions{Size: size}
}
