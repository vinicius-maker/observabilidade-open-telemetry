package infraestruct

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/vinicius-maker/observabilidade-open-telemetry/internal/entity"
	"io"
	"log"
	"net/http"
)

var ErrRequestZipCode = errors.New("an error occurred while processing your request of zipcode")

type ViaCepStruct struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ServiceViaCep struct {
	httpTracer *HttpTracer
}

func NewServiceViaCep(tracer *HttpTracer) *ServiceViaCep {
	return &ServiceViaCep{
		httpTracer: tracer,
	}
}

func (s *ServiceViaCep) SearchCep(ctx context.Context, cepCode *entity.CepCode) (string, error) {
	res, err := s.httpTracer.Get(ctx, "http://viacep.com.br/ws/"+cepCode.CepCode+"/json/", "call viacep service")
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Error: unexpected status code %d from ViaCEP API response", res.StatusCode)
		return "", ErrRequestZipCode
	}

	if err != nil {
		log.Printf("Error: %v", err)
		return "", ErrRequestZipCode
	}

	body, err := io.ReadAll(res.Body)

	var data ViaCepStruct
	err = json.Unmarshal(body, &data)

	return data.Localidade, err
}
