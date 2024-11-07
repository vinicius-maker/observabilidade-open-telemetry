package service

import "context"

type WeatherService interface {
	DiscoverWeather(ctx context.Context, location string) (float64, error)
}
