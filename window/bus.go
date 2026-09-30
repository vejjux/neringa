package window

import (
	"gioui.org/layout"
	"keltas/schedule"
)

func Bus(b schedule.Bus, h *Home) func(gtx layout.Context) {
	return func(gtx layout.Context) {
		stop := Stops[h.Selected]
		Schedule(schedule.Ferry{
			Title: "KAUTRA",
			Schedules: []schedule.Schedule{
				{Title: "Iš Smiltynės", Table: b.Back},
				{Title: "Iš " + stop, Table: b.There[stop]},
			},
		})(gtx)
	}
}
