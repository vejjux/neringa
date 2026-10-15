package schedule

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const (
	Nida     = "Nida"
	Smiltyne = "Smiltynė"
)

// BusPage lists the monthly timetable published by the Neringa municipality.
const BusPage = "/gyventojui/autobusu-tvarkarastis-ir-bilietu-kainos/63"

var (
	timetableLink = regexp.MustCompile(`href="(/data/public/uploads/[^"]*autobusu-grafikas[^"]*\.docx)"`)
	departureTime = regexp.MustCompile(`^(\d{1,2}\.\d{2})(d\.d\.)?`)
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

func FetchBus(origin string) (bus Bus, err error) {
	page, err := get(origin + BusPage)
	if err != nil {
		return
	}
	m := timetableLink.FindSubmatch(page)
	if m == nil {
		return bus, errors.New("bus timetable not found")
	}
	doc, err := get(origin + string(m[1]))
	if err != nil {
		return
	}
	tables, err := docxTables(doc)
	if err != nil {
		return
	}

	bus.ToSmiltyne = make(map[string][]Time)
	bus.ToNida = make(map[string][]Time)
	for _, rows := range tables {
		var stops []string
		var columns [][]string
		for _, cells := range rows {
			if len(cells) < 2 {
				continue
			}
			if stops == nil {
				if strings.HasPrefix(strings.TrimSpace(cells[1]), "Stotelė") {
					for _, c := range cells[1:] {
						stops = append(stops, StopName(c))
					}
					columns = make([][]string, len(stops))
				}
				continue
			}
			for i, c := range cells[1:] {
				if t := departure(c); i < len(columns) && t != "" {
					columns[i] = append(columns[i], t)
				}
			}
		}
		if len(stops) == 0 {
			continue
		}

		direction := bus.ToNida
		if stops[0] == Nida {
			direction = bus.ToSmiltyne
		}
		for i, s := range stops {
			direction[s] = group(columns[i])
		}
	}
	if len(bus.ToSmiltyne) == 0 && len(bus.ToNida) == 0 {
		err = errors.New("no bus departures found")
	}

	return
}

func get(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, r.Status)
	}
	return io.ReadAll(r.Body)
}

// docxTables returns the text of every table cell in a .docx document,
// as tables of rows of cells.
func docxTables(doc []byte) (tables [][][]string, err error) {
	z, err := zip.NewReader(bytes.NewReader(doc), int64(len(doc)))
	if err != nil {
		return
	}
	f, err := z.Open("word/document.xml")
	if err != nil {
		return
	}
	defer f.Close()

	var rows [][]string
	var text bool
	d := xml.NewDecoder(f)
	for {
		tok, err := d.Token()
		if err == io.EOF {
			return tables, nil
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "tbl":
				rows = nil
			case "tr":
				rows = append(rows, nil)
			case "tc":
				if n := len(rows); n > 0 {
					rows[n-1] = append(rows[n-1], "")
				}
			case "t":
				text = true
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "tbl":
				tables = append(tables, rows)
			case "t":
				text = false
			}
		case xml.CharData:
			if n := len(rows); text && n > 0 && len(rows[n-1]) > 0 {
				rows[n-1][len(rows[n-1])-1] += string(t)
			}
		}
	}
}

// departure normalizes a timetable cell, e.g. "6.00 d. d." into "6.00D",
// where D marks departures on workdays only. Anything after the time, such
// as a footnote, is dropped; cells without a time yield "".
func departure(s string) string {
	s = strings.Join(strings.Fields(s), "")
	m := departureTime.FindStringSubmatch(s)
	if m == nil {
		return ""
	}
	// "d. d. – darbo dienomis" after the last time is the legend, not a mark.
	if m[2] != "" && !strings.HasPrefix(s[len(m[0]):], "–") {
		return m[1] + "D"
	}
	return m[1]
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
