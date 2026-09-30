package schedule

import (
	"fmt"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"keltas/schedule"
	"keltas/widgets"
)

func Times(th *material.Theme, times []schedule.Time) layout.Widget {
	labels := make([]layout.FlexChild, 0, len(times))
	for i, t := range times {
		hour := widgets.NewLabel(th, 20, fmt.Sprintf(" %s: ", t.Hour))
		hour.Color = ColorTableText
		hour.Font.Weight = font.Bold
		hour.Alignment = text.Start

		minutes := widgets.NewLabel(th, 16, t.Minutes)
		minutes.Color = ColorTableText
		minutes.Font.Weight = font.Bold
		minutes.Alignment = text.Start

		if i%2 == 1 {
			hour.Background, minutes.Background = ColorStripe, ColorStripe
		}

		labels = append(labels, layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Start,
				}.Layout(gtx, layout.Rigid(hour.Layout), layout.Rigid(minutes.Layout))
			},
		))
	}

	return func(gtx layout.Context) (dims layout.Dimensions) {
		flexWidget := func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Start,
			}.Layout(gtx, labels...)
		}

		widget.Border{
			Color:        ColorBorder,
			CornerRadius: 0,
			Width:        3,
		}.Layout(gtx, flexWidget)

		return
	}
}
