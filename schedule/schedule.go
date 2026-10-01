package schedule

import "strconv"

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

func (f Ferry) Upcoming(hour, n int) Ferry {
	schedules := make([]Schedule, len(f.Schedules))
	for i, s := range f.Schedules {
		schedules[i] = Schedule{Title: s.Title, Table: upcoming(s.Table, hour, n)}
	}
	f.Schedules = schedules
	return f
}

func (b Bus) Upcoming(hour, n int) Bus {
	return Bus{ToSmiltyne: upcomingAll(b.ToSmiltyne, hour, n), ToNida: upcomingAll(b.ToNida, hour, n)}
}

func upcomingAll(stops map[string][]Time, hour, n int) map[string][]Time {
	result := make(map[string][]Time, len(stops))
	for k, v := range stops {
		result[k] = upcoming(v, hour, n)
	}
	return result
}

func upcoming(table []Time, hour, n int) []Time {
	i := 0
	for ; i < len(table); i++ {
		if h, err := strconv.Atoi(table[i].Hour); err != nil || h >= hour {
			break
		}
	}
	rotated := append(append([]Time{}, table[i:]...), table[:i]...)
	return rotated[:min(n, len(rotated))]
}
