package request

type DateHistogram struct {
	Field         string `json:"field"`
	FixedInterval string `json:"fixed_interval"`
	Format        string `json:"format"`
	Timezone      string `json:"time_zone"`
}

type ByApps struct {
	TermsInclude    `json:"terms"`
	AggsMyBandwidth `json:"aggs"`
}

type ByOthers struct {
	TermsExclude    `json:"terms"`
	AggsMyBandwidth `json:"aggs"`
}

type TermsInclude struct {
	Field   string `json:"field"`
	Include string `json:"include"`
}

type TermsExclude struct {
	Field   string `json:"field"`
	Include string `json:"exclude"`
}

type Aggs struct {
	ByApps   `json:"by_apps"`
	ByOthers `json:"by_others"`
}

type AggsMyBandwidth struct {
	MyBandwidth `json:"my_Bandwidth"`
}

type MyBandwidth struct {
	Sum `json:"sum"`
}

type Sum struct {
	Field string `json:"field"`
}

type MyDateHistogram struct {
	DateHistogram DateHistogram `json:"date_histogram"`
	Aggs          `json:"aggs"`
}

type Aggregations struct {
	MyDateHistogram `json:"myDateHistogram"`
}

/** ****************QueryForTemporalEvolution Objects********************* **/

type TimeStamp struct {
	Format string `json:"format"`
	Gte    string `json:"gte"`
	Lte    string `json:"lte"`
}

type Range struct {
	TimeStamp `json:"@timestamp"`
}

type RangeObj struct {
	Range `json:"range"`
}

type Bool struct {
	Filter []RangeObj `json:"filter"`
}

type Query struct {
	Bool `json:"bool"`
}

type TemporalEvolutionRequest struct {
	Size         int `json:"size"`
	Query        `json:"query"`
	Aggregations `json:"aggregations"`
}

/**   Temporal Selection Request **/
type TemporalSelectionRequest struct {
	Size                  int `json:"size"`
	Query                 `json:"query"`
	AggregationsSelection `json:"aggregations"`
}

type AggregationsSelection struct {
	MyDateHistogramSelection `json:"myDateHistogram"`
}

type MyDateHistogramSelection struct {
	DateHistogram   DateHistogram `json:"date_histogram"`
	AggsMyBandwidth `json:"aggs"`
}
