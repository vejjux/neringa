package window

import (
	"fmt"
	"image/color"
	"math"
	"strings"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"keltas/schedule"
	"keltas/weather"
	ui "keltas/window/schedule"
)

const tileDepartures = 2

// Tile is a home page box that opens Page when clicked.
type Tile struct {
	Page    int
	Content func(gtx layout.Context)
}

type departures struct {
	title string
	times []string
}

// FerryTile shows the next departures from each side, Smiltynė first.
func FerryTile(f schedule.Ferry) func(gtx layout.Context) {
	return func(gtx layout.Context) {
		now := time.Now()
		rows := make([]departures, 0, len(f.Schedules))
		for i := len(f.Schedules) - 1; i >= 0; i-- {
			s := f.Schedules[i]
			rows = append(rows, departures{s.Title, schedule.Next(s.Table, now, tileDepartures)})
		}
		layoutDepartures(gtx, rows)
	}
}

// BusTile shows the next departures from the home stop in each direction.
func BusTile(b schedule.Bus, h *Home) func(gtx layout.Context) {
	return func(gtx layout.Context) {
		now := time.Now()
		stop := Stops[h.Selected]
		var rows []departures
		for _, s := range busSchedules(b, stop) {
			title := "→" + strings.TrimPrefix(s.Title, stop+" →")
			rows = append(rows, departures{title, schedule.Next(s.Table, now, tileDepartures)})
		}
		layoutDepartures(gtx, rows)
	}
}

// WeatherTile shows the current hour's forecast for the home stop.
func WeatherTile(forecasts map[string][]weather.Hour, h *Home) func(gtx layout.Context) {
	return func(gtx layout.Context) {
		place := StopPlaces[Stops[h.Selected]]
		hours := forecasts[place.Code]
		if len(hours) == 0 {
			tileLabel(gtx, 14, font.Normal, ui.ColorOnHeader, "Nėra duomenų")
			return
		}
		f := hours[0]
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return tileLabel(gtx, 12, font.Normal, ui.ColorOnHeader, place.Name+" "+f.Time.Local().Format("15:04"))
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return tileLabel(gtx, 36, font.Bold, ui.ColorTableText, fmt.Sprintf("%d°", int(math.Round(f.Temperature))))
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				wind := fmt.Sprintf("%s %d m/s", compass[int(math.Round(f.WindDirection/45))%8], int(math.Round(f.WindSpeed)))
				return tileLabel(gtx, 16, font.Bold, ui.ColorTableText, wind)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return tileLabel(gtx, 16, font.Bold, ui.ColorTableText, fmt.Sprintf("%.1f mm", f.Precipitation))
			}),
		)
	}
}

// TileError shows err in a size that fits a tile.
func TileError(err error) func(gtx layout.Context) {
	return func(gtx layout.Context) {
		th := newTheme()
		body := material.Caption(th, err.Error())
		body.Color = ui.ColorError
		body.MaxLines = 5
		body.Layout(gtx)
	}
}

func layoutDepartures(gtx layout.Context, rows []departures) {
	children := make([]layout.FlexChild, 0, 2*len(rows))
	for _, r := range rows {
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return tileLabel(gtx, 12, font.Normal, ui.ColorOnHeader, r.title)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: 6}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return tileLabel(gtx, 20, font.Bold, ui.ColorTableText, strings.Join(r.times, "  "))
				})
			}),
		)
	}
	layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func tileLabel(gtx layout.Context, size unit.Sp, weight font.Weight, color color.NRGBA, txt string) layout.Dimensions {
	l := material.Label(newTheme(), size, txt)
	l.Color = color
	l.Font.Weight = weight
	l.MaxLines = 1
	return l.Layout(gtx)
}
