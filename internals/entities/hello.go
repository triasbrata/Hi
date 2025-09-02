package entities

type FetchCurrentWeatherParam struct {
	Latitude  float32
	Longitude float32
}
type FetchCurrentWeatherRes struct {
	Temperature float32
	Condition   string
}
