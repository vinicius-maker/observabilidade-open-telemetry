package service

import (
	"context"
	"github.com/vinicius-maker/observabilidade-open-telemetry/internal/entity"
)

type ViaCepService interface {
	SearchCep(ctx context.Context, cepCode *entity.CepCode) (string, error)
}
