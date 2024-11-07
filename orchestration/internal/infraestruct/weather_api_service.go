package infraestruct

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

type WeatherReturn struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherApiService struct {
	apiKey     string
	httpTracer *HttpTracer
}

func NewWeatherServiceWeatherApi(apiKey string, tracer *HttpTracer) *WeatherApiService {
	return &WeatherApiService{
		apiKey:     apiKey,
		httpTracer: tracer,
	}
}

func (w *WeatherApiService) DiscoverWeather(ctx context.Context, location string) (float64, error) {
	encodedLocation := url.QueryEscape(location)
	urlStr := "http://api.weatherapi.com/v1/current.json?key=" + w.apiKey + "&q=" + encodedLocation

	res, err := w.httpTracer.Get(ctx, urlStr, "call weatherapi service")
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("error status code %d from Wheater API", res.StatusCode)
		return 0.00, errors.New("error occurred while processing your request of Wheater API response")
	}

	if err != nil {
		return 0.00, err
	}

	body, err := io.ReadAll(res.Body)

	var WeatherReturn *WeatherReturn

	errJson := json.Unmarshal(body, &WeatherReturn)

	if errJson != nil {
		return 0.00, errJson
	}

	return WeatherReturn.Current.TempC, nil
}
