package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/vinicius-maker/observabilidade-open-telemetry/controller"
	"github.com/vinicius-maker/observabilidade-open-telemetry/internal/infraestruct"
	"github.com/vinicius-maker/observabilidade-open-telemetry/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("not found environment variable API_KEY")
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := infraestruct.Provider()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	httpTracer := infraestruct.NewHttpTracer("orchestration")

	viaCepService := infraestruct.NewServiceViaCep(httpTracer)
	weatherService := infraestruct.NewWeatherServiceWeatherApi(apiKey, httpTracer)

	app := usecase.NewDiscoverWeatherByLocation(viaCepService, weatherService)

	weatherController := controller.NewWeatherController(app)

	http.HandleFunc("/discover-temperature", weatherController.Handle)

	http.ListenAndServe(":8080", nil)
}
