package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexchebotarsky/heatpump-api/client"
	"github.com/alexchebotarsky/heatpump-api/model/heatpump"
)

type TemperatureAndHumidityFetcher interface {
	FetchTemperatureAndHumidity() (temperature float64, humidity float64, err error)
}

func GetTemperatureAndHumidity(fetcher TemperatureAndHumidityFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temperature, humidity, err := fetcher.FetchTemperatureAndHumidity()
		if err != nil {
			switch err.(type) {
			case *client.ErrNotFound:
				log.Printf("Temperature and humidity not found: %v", err)
			default:
				HandleError(w, fmt.Errorf("error fetching temperature and humidity: %v", err), http.StatusInternalServerError, true)
				return
			}
		}

		temperatureReading := heatpump.TemperatureReading{
			Temperature: temperature,
			Humidity:    humidity,
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(temperatureReading)
		handleWritingErr(err)
	}
}
