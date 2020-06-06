package domain

type TemporalEvolutionResponse struct {
	Aggregations MyAggregations `json:"aggregations"`
	Shards       `json:"_shards"`
	TimedOut     bool `json:"timed_out"`
	Took         int  `json:"took"`
}

//type Hits struct {
//	Hits     []string `json:"hits"`
//	MaxScore string   `json:"max_score"`
//	Total    string   `json:"total"`
//}

type MyTotal struct {
	Relation string `json:"relation"`
	Value    string `json:"value"`
}

type Shards struct {
	Failed     int `json:"failed"`
	Skipped    int `json:"skipped"`
	Successful int `json:"successful"`
	MyTotal    int `json:"total"`
}

type MyAggregations struct {
	MyDateHistogram `json:"myDateHistogram"`
}

type MyDateHistogram struct {
	Buckets []DateHistogramBucket `json:"buckets"`
}

type DateHistogramBucket struct {
	ByApps      Apps   `json:"by_apps"`
	ByOthers    Apps   `json:"by_others"`
	DocCount    int    `json:"doc_count"`
	Key         int64  `json:"key"`
	KeyAsString string `json:"key_as_string"`
}

type Apps struct {
	Buckets                 []Bucket `json:"buckets"`
	DocCountErrorUpperBound int      `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int      `json:"sum_other_doc_count"`
}

type Bucket struct {
	MyBandWidth `json:"my_Bandwidth"`
	DocCount    int    `json:"doc_count"`
	Key         string `json:"key"`
}

type MyBandWidth struct {
	Value float32 `json:"value"`
}

/** TemporalSelectionResponse  **/
type TemporalSelectionResponse struct {
	AggregationsSelection `json:"aggregations"`
	Shards                `json:"_shards"`
	TimedOut              bool `json:"timed_out"`
	Took                  int  `json:"took"`
}

type DateHistogramBucketSelection struct {
	DocCount    int    `json:"doc_count"`
	Key         int64  `json:"key"`
	KeyAsString string `json:"key_as_string"`
	MyBandWidth `json:"my_Bandwidth"`
}

type AggregationsSelection struct {
	MyDateHistogramSelection `json:"myDateHistogram"`
}

type MyDateHistogramSelection struct {
	Buckets []DateHistogramBucketSelection `json:"buckets"`
}
