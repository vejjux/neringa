package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var Places = []Place{Nida, Preila, Pervalka, Juodkrante, Smiltyne}

type Hour struct {
	Time          time.Time `json:"-"`
	TimeUTC       string    `json:"forecastTimeUtc"`
	Temperature   float64   `json:"airTemperature"`
	WindSpeed     float64   `json:"windSpeed"`
	WindDirection float64   `json:"windDirection"`
	Precipitation float64   `json:"totalPrecipitation"`
}

func Fetch(base string, places ...Place) (map[string][]Hour, error) {
	forecasts := make(map[string][]Hour)
	for _, p := range places {
		hours, err := fetch(base + "/v1/places/" + p.Code + "/forecasts/long-term")
		if err != nil {
			return nil, err
		}
		forecasts[p.Code] = hours
	}
	return forecasts, nil
}

func fetch(url string) (hours []Hour, err error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, r.Status)
	}

	var forecast struct {
		ForecastTimestamps []Hour `json:"forecastTimestamps"`
	}
	if err = json.NewDecoder(r.Body).Decode(&forecast); err != nil {
		return nil, err
	}

	from := time.Now().Truncate(time.Hour)
	to := from.Add(24 * time.Hour)
	for _, h := range forecast.ForecastTimestamps {
		if h.Time, err = time.Parse(time.DateTime, h.TimeUTC); err != nil {
			return nil, err
		}
		if !h.Time.Before(from) && h.Time.Before(to) {
			hours = append(hours, h)
		}
	}
	return hours, nil
}
