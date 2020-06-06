package response

type EvolutionObject struct {
	BorderWidth     int       `json:"borderWidth"`
	BorderColor     string    `json:"borderColor"`
	BackgroundColor string    `json:"backgroundColor"`
	Label           string    `json:"label"`
	Data            []float32 `json:"data"`
	BarData         float32   `json:"bar_data"`
	Count           int       `json:"count"`
}

type UiData struct {
	Labels   []string          `json:"labels"`
	DataSets []EvolutionObject `json:"datasets"`
}

type UiBarObject struct {
	Labels   []string         `json:"labels"`
	DataSets []FinalBarObject `json:"datasets"`
}

type FinalBarObject struct {
	BarThickness    int       `json:"barThickness"`
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	BackgroundColor []string  `json:"backgroundColor"`
}

type EvolutionResponse struct {
	UiData      `json:"ui_data"`
	UiBarObject `json:"ui_bar_object"`
}
