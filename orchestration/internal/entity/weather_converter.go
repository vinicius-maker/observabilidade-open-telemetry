package entity

type WeatherConverter struct {
	Celsius float64
}

func NewWeatherConverter(celsius float64) *WeatherConverter {
	return &WeatherConverter{Celsius: celsius}
}

func (w *WeatherConverter) ToFahrenheit() float64 {
	return (w.Celsius * 1.8) + 32
}

func (w *WeatherConverter) ToKelvin() float64 {
	return w.Celsius + 273
}
