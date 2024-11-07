package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyZipCode_ShouldReturnError(t *testing.T) {
	cep, err := NewCepCode("")

	assert.Equal(t, "invalid zipcode", err.Error())
	assert.Nil(t, cep)
}

func TestShortZipCode_ShouldReturnError(t *testing.T) {
	cep, err := NewCepCode("12345")

	assert.Equal(t, "invalid zipcode", err.Error())
	assert.Nil(t, cep)
}

func TestValidZipCode_ShouldCreateZipCode(t *testing.T) {
	cepCode, err := NewCepCode("98765432")

	assert.Equal(t, "98765432", cepCode.CepCode)
	assert.Nil(t, err)
}
