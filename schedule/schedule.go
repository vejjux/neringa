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
	there := make(map[string][]Time, len(b.There))
	for k, v := range b.There {
		there[k] = upcoming(v, hour, n)
	}
	return Bus{There: there, Back: upcoming(b.Back, hour, n)}
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
