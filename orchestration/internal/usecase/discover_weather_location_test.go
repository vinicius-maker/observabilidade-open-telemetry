package usecase

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/vinicius-maker/observabilidade-open-telemetry/internal/infraestruct"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func LoadUseCase() *DiscoverWeatherByLocation {
	envFile := filepath.Join("..", "..", ".env")

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	httpTracer := infraestruct.NewHttpTracer("orchestration")

	zipCodeService := infraestruct.NewServiceViaCep(httpTracer)
	weatherService := infraestruct.NewWeatherServiceWeatherApi(apiKey, httpTracer)

	return NewDiscoverWeatherByLocation(zipCodeService, weatherService)
}

func TestShortZipCode_ShouldReturnErrorMessage(t *testing.T) {
	useCase := LoadUseCase()

	var dto DiscoverWeatherByLocationDTO
	dto.CepCode = "56789"

	ctx := context.Background()
	_, err := useCase.Execute(ctx, dto)

	assert.EqualError(t, err, "invalid zipcode")
}

func TestNonExistentZipCode_ShouldReturnNotFoundError(t *testing.T) {
	useCase := LoadUseCase()

	var dto DiscoverWeatherByLocationDTO
	dto.CepCode = "11111111"

	ctx := context.Background()
	_, err := useCase.Execute(ctx, dto)

	assert.EqualError(t, err, "can not find zipcode")
}

func TestValidZipCode_ShouldReturnWeatherDetails(t *testing.T) {
	useCase := LoadUseCase()

	var dto DiscoverWeatherByLocationDTO
	dto.CepCode = "07115260"

	ctx := context.Background()
	output, err := useCase.Execute(ctx, dto)

	assert.Nil(t, err)
	assert.IsType(t, float64(0), output.TempC)
	assert.IsType(t, float64(0), output.TempF)
	assert.IsType(t, float64(0), output.TempK)
}
