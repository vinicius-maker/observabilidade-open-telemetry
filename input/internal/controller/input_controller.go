package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vinicius-maker/observabilidade-open-telemetry/internal/infraestruct"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"
)

type RequestParams struct {
	CEP string `json:"cep"`
}

func isValidCEP(cepParam string) bool {
	var validCEP = regexp.MustCompile(`^\d{5}-?\d{3}$`)
	return validCEP.MatchString(cepParam)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func Handle(w http.ResponseWriter, r *http.Request) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var reqParams RequestParams
	fmt.Println(reqParams)
	if err := json.NewDecoder(r.Body).Decode(&reqParams); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if !isValidCEP(reqParams.CEP) {
		http.Error(w, "Invalid CEP", http.StatusUnprocessableEntity)
		return
	}

	shutdown, err := infraestruct.Provider()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("microservice-tracer")

	carrier := propagation.HeaderCarrier(r.Header)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := tracer.Start(ctx, "call orchestration")
	defer span.End()

	time.Sleep(time.Second * 3)

	url := os.Getenv("ORCHESTRATION_URL") + "discover-temperature?cep=" + reqParams.CEP

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "GET", url, nil)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Error: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Error processing the request")
		return
	}

	if err != nil {
		log.Printf("Error: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Error processing the request")
		return
	}

	body, err := io.ReadAll(res.Body)

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Printf("Error: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Error processing the request")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(data)
}
