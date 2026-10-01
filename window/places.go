package window

import (
	"keltas/schedule"
	"keltas/weather"
)

var StopPlaces = map[string]weather.Place{
	schedule.Nida:                   weather.Nida,
	"G. D. Kuverto plento sankryža": weather.Nida,
	"T. Mano muziejus":              weather.Nida,
	"Preila":                        weather.Preila,
	"Preila prie plento":            weather.Preila,
	"Pervalka":                      weather.Pervalka,
	"Pervalka prie plento":          weather.Pervalka,
	"Žvejų kaimelis":                weather.Juodkrante,
	"Raganų kalnas":                 weather.Juodkrante,
	"Juodkrantė":                    weather.Juodkrante,
	"Gintaro įlanka":                weather.Juodkrante,
	"Alksnynė":                      weather.Smiltyne,
}
