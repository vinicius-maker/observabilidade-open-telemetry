package usecase

import (
	"context"
	"errors"
	"github.com/vinicius-maker/observabilidade-open-telemetry/internal/entity"
	"github.com/vinicius-maker/observabilidade-open-telemetry/internal/service"
)

var ErrCepCodeNotFound = errors.New("can not find zipcode")

type DiscoverWeatherByLocationDTO struct {
	CepCode string `json:"zip_code"`
}

type DiscoverWeatherByLocationOutputDTO struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type DiscoverWeatherByLocation struct {
	viaCepService  service.ViaCepService
	weatherService service.WeatherService
}

func NewDiscoverWeatherByLocation(viaCepService service.ViaCepService, weatherService service.WeatherService) *DiscoverWeatherByLocation {
	return &DiscoverWeatherByLocation{
		viaCepService:  viaCepService,
		weatherService: weatherService,
	}
}

func (d *DiscoverWeatherByLocation) Execute(ctx context.Context, inputDTO DiscoverWeatherByLocationDTO) (DiscoverWeatherByLocationOutputDTO, error) {
	outputDTO := DiscoverWeatherByLocationOutputDTO{}

	cepCode, err := entity.NewCepCode(inputDTO.CepCode)

	if err != nil {
		return outputDTO, err
	}

	location, err := d.viaCepService.SearchCep(ctx, cepCode)
	if err != nil {
		return outputDTO, err
	}

	if location == "" {
		return outputDTO, ErrCepCodeNotFound
	}

	weather, err := d.weatherService.DiscoverWeather(ctx, location)
	if err != nil {
		return outputDTO, err
	}

	converter := entity.WeatherConverter{
		Celsius: weather,
	}

	outputDTO.City = location
	outputDTO.TempC = weather
	outputDTO.TempF = converter.ToFahrenheit()
	outputDTO.TempK = converter.ToKelvin()

	return outputDTO, err
}
