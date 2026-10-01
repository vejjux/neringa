package window

import (
	"fmt"
	"gioui.org/layout"
	"keltas/schedule"
	"keltas/widgets"
	ui "keltas/window/schedule"
)

// Schedule lays out one or two schedules side by side, the last one on the left.
func Schedule(f schedule.Ferry) func(gtx layout.Context) {
	if n := len(f.Schedules); n < 1 || n > 2 {
		return Error(fmt.Errorf("no schedules found"))
	}

	return func(gtx layout.Context) {
		th := newTheme()
		column := func(s schedule.Schedule) layout.Widget {
			return func(g layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Start,
				}.Layout(g, layout.Rigid(ui.Direction(th, s.Title)), layout.Rigid(ui.Times(th, s.Table)))
			}
		}

		body := column(f.Schedules[0])
		if len(f.Schedules) == 2 {
			body = widgets.Split{
				Left:  column(f.Schedules[1]),
				Right: column(f.Schedules[0]),
			}.Layout
		}

		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Start,
		}.Layout(gtx, layout.Rigid(ui.Title(th, f.Title)), layout.Flexed(100, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = gtx.Constraints.Max
			return body(gtx)
		}))
	}
}

// Stack lays out sections vertically, each sized by its weight.
func Stack(weights []float32, sections ...func(layout.Context)) func(layout.Context) {
	return func(gtx layout.Context) {
		children := make([]layout.FlexChild, len(sections))
		for i, section := range sections {
			children[i] = layout.Flexed(weights[i], func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: 4}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = gtx.Constraints.Max
					section(gtx)
					return layout.Dimensions{Size: gtx.Constraints.Max}
				})
			})
		}
		layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	}
}
