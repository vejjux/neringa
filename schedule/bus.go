package schedule

import (
	"strings"

	"github.com/gocolly/colly"
)

type Bus struct {
	There map[string][]Time
	Back  []Time
}

func FetchBus(url string) (bus Bus, err error) {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) { r.ResponseCharacterEncoding = "windows-1257" })
	bus.There = make(map[string][]Time)

	tables := 0
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.DOM.Find("th").Length() == 0 {
			return
		}
		tables++
		if tables > 2 {
			return
		}

		var stops []string
		e.ForEach("th", func(_ int, e *colly.HTMLElement) {
			stops = append(stops, stopName(e.Text))
		})
		columns := make([][]string, len(stops))
		e.ForEach("tr", func(_ int, e *colly.HTMLElement) {
			e.ForEach("td", func(i int, e *colly.HTMLElement) {
				if t := strings.Join(strings.Fields(e.Text), ""); i < len(columns) && t != "" {
					columns[i] = append(columns[i], t)
				}
			})
		})

		if tables == 1 {
			for i, s := range stops {
				bus.There[s] = group(columns[i])
			}
		} else if len(columns) > 0 {
			bus.Back = group(columns[0])
		}
	})

	err = c.Visit(url)

	return
}

func stopName(s string) string {
	s = strings.TrimPrefix(strings.TrimSpace(s), "Stotelė")
	return strings.Join(strings.Fields(strings.Trim(s, " „“")), " ")
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
