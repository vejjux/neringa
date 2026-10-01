package schedule

import (
	"strings"

	"github.com/gocolly/colly"
)

const (
	Nida     = "Nida"
	Smiltyne = "Smiltynė"
)

var stopNames = map[string]string{
	"Nidos gyvenvietės autobusų stotis": Nida,
	"Preilos gv.":                       "Preila",
	"Preilos gv. prie plento":           "Preila prie plento",
	"Pervalkos gv.":                     "Pervalka",
	"Pervalkos gv. prie plento":         "Pervalka prie plento",
	"Juodkrantės gv.":                   "Juodkrantė",
}

// Bus holds departure times per stop for each direction.
type Bus struct {
	ToSmiltyne map[string][]Time
	ToNida     map[string][]Time
}

func FetchBus(url string) (bus Bus, err error) {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) { r.ResponseCharacterEncoding = "windows-1257" })
	bus.ToSmiltyne = make(map[string][]Time)
	bus.ToNida = make(map[string][]Time)

	c.OnHTML("table", func(e *colly.HTMLElement) {
		var stops []string
		e.ForEach("th", func(_ int, e *colly.HTMLElement) {
			stops = append(stops, StopName(e.Text))
		})
		if len(stops) == 0 {
			return
		}
		columns := make([][]string, len(stops))
		e.ForEach("tr", func(_ int, e *colly.HTMLElement) {
			e.ForEach("td", func(i int, e *colly.HTMLElement) {
				if t := strings.Join(strings.Fields(e.Text), ""); i < len(columns) && t != "" {
					columns[i] = append(columns[i], t)
				}
			})
		})

		direction := bus.ToNida
		if stops[0] == Nida {
			direction = bus.ToSmiltyne
		}
		for i, s := range stops {
			direction[s] = group(columns[i])
		}
	})

	err = c.Visit(url)

	return
}

// StopName normalizes a stop title from the timetable into its display name.
func StopName(s string) string {
	s = strings.TrimPrefix(strings.TrimSpace(s), "Stotelė")
	s = strings.Join(strings.Fields(strings.Trim(s, " „“")), " ")
	if n, ok := stopNames[s]; ok {
		return n
	}
	return s
}

func group(times []string) (table []Time) {
	for _, t := range times {
		h, m, _ := strings.Cut(t, ".")
		if n := len(table); n > 0 && table[n-1].Hour == h {
			table[n-1].Minutes += " " + m
			continue
		}
		table = append(table, Time{Hour: h, Minutes: m})
	}
	return
}
