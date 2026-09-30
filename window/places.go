package window

import "keltas/weather"

var StopPlaces = map[string]weather.Place{
	"Nidos gyvenvietės autobusų stotis": weather.Nida,
	"G. D. Kuverto plento sankryža":     weather.Nida,
	"T. Mano muziejus":                  weather.Nida,
	"Preilos gv.":                       weather.Preila,
	"Preilos gv. prie plento":           weather.Preila,
	"Pervalkos gv.":                     weather.Pervalka,
	"Pervalkos gv. prie plento":         weather.Pervalka,
	"Žvejų kaimelis":                    weather.Juodkrante,
	"Raganų kalnas":                     weather.Juodkrante,
	"Juodkrantės gv.":                   weather.Juodkrante,
	"Gintaro įlanka":                    weather.Juodkrante,
	"Alksnynė":                          weather.Smiltyne,
}
