package entity

import (
	"errors"
	"regexp"
)

var ErrInvalidCepCode = errors.New("invalid zipcode")

type CepCode struct {
	CepCode string
}

func NewCepCode(cepParam string) (*CepCode, error) {
	cepCode := &CepCode{
		CepCode: cepParam,
	}

	err := cepCode.isValidCEP()
	if err != nil {
		return nil, err
	}

	return cepCode, nil
}

func (c *CepCode) isValidCEP() error {
	// "12345678" or "12345-678"
	validCEP := regexp.MustCompile(`^\d{5}-?\d{3}$`)
	if !validCEP.MatchString(c.CepCode) {
		return ErrInvalidCepCode
	}

	return nil
}
