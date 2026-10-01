package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	Hour    string
	Minutes string
}

type Schedule struct {
	Title string
	Table []Time
}

type Ferry struct {
	Title     string
	Schedules []Schedule
}

// Next returns the next n departures in table at or after now as "15:04",
// wrapping around to the start of the day.
func Next(table []Time, now time.Time, n int) []string {
	var all []int
	for _, t := range table {
		h, err := strconv.Atoi(t.Hour)
		if err != nil {
			continue
		}
		for _, m := range strings.Fields(t.Minutes) {
			if m, err := strconv.Atoi(strings.TrimRight(m, "*D")); err == nil {
				all = append(all, h*60+m)
			}
		}
	}
	current := now.Hour()*60 + now.Minute()
	i := 0
	for i < len(all) && all[i] < current {
		i++
	}
	var next []string
	for k := 0; k < min(n, len(all)); k++ {
		m := all[(i+k)%len(all)]
		next = append(next, fmt.Sprintf("%02d:%02d", m/60, m%60))
	}
	return next
}
