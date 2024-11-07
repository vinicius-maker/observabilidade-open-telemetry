package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCelsiusToFahrenheit_ShouldReturnCorrectValue(t *testing.T) {
	converter := NewWeatherConverter(32)

	assert.Equal(t, 89.6, converter.ToFahrenheit())
}

func TestCelsiusToKelvin_ShouldReturnCorrectValue(t *testing.T) {
	converter := NewWeatherConverter(32)

	assert.Equal(t, 305.0, converter.ToKelvin())
}
