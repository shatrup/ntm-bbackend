package request

type TemporalRequest struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Interval string `json:"interval"`
}
