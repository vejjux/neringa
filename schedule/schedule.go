package schedule

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

func (f Ferry) First(n int) Ferry {
	schedules := make([]Schedule, len(f.Schedules))
	for i, s := range f.Schedules {
		schedules[i] = Schedule{Title: s.Title, Table: s.Table[:min(n, len(s.Table))]}
	}
	f.Schedules = schedules
	return f
}

func (b Bus) First(n int) Bus {
	there := make(map[string][]Time, len(b.There))
	for k, v := range b.There {
		there[k] = v[:min(n, len(v))]
	}
	return Bus{There: there, Back: b.Back[:min(n, len(b.Back))]}
}
