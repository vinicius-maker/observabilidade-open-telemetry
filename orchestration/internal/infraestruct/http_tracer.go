package infraestruct

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"net/http"
)

type HttpTracer struct {
	request     *http.Request
	ctx         context.Context
	serviceName string
}

func NewHttpTracer(serviceName string) *HttpTracer {
	return &HttpTracer{
		serviceName: serviceName,
	}
}

func (h *HttpTracer) Get(ctx context.Context, url, spanName string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	tracer := otel.Tracer("microservice-tracer")

	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	res, err := http.DefaultClient.Do(req)

	return res, err
}
