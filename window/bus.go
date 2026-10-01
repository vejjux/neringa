package window

import (
	"gioui.org/layout"
	"keltas/schedule"
)

const busTitle = "KAUTRA"

// Bus shows departures from the home stop in each direction it has.
func Bus(b schedule.Bus, h *Home) func(gtx layout.Context) {
	return func(gtx layout.Context) {
		stop := Stops[h.Selected]
		var schedules []schedule.Schedule
		if t := b.ToSmiltyne[stop]; stop != schedule.Smiltyne && len(t) > 0 {
			schedules = append(schedules, schedule.Schedule{Title: stop + " → " + schedule.Smiltyne, Table: t})
		}
		if t := b.ToNida[stop]; stop != schedule.Nida && len(t) > 0 {
			schedules = append(schedules, schedule.Schedule{Title: stop + " → " + schedule.Nida, Table: t})
		}
		Schedule(schedule.Ferry{Title: busTitle, Schedules: schedules})(gtx)
	}
}

// BusFromSmiltyne shows departures from Smiltynė towards Nida.
func BusFromSmiltyne(b schedule.Bus) func(gtx layout.Context) {
	return Schedule(schedule.Ferry{
		Title: busTitle,
		Schedules: []schedule.Schedule{
			{Title: schedule.Smiltyne + " → " + schedule.Nida, Table: b.ToNida[schedule.Smiltyne]},
		},
	})
}
