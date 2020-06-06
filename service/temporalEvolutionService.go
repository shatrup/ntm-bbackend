package service

import (
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"ntm-backend/domain"
	"ntm-backend/dto/request"
	"ntm-backend/dto/response"
	"ntm-backend/repository"
	"time"
)

type TemporalEvolutionService struct {
	Logger     *logrus.Logger
	Repository repository.ESRepository
}

func (ts TemporalEvolutionService) BuildTemporalEvolutionData(req request.TemporalRequest) (response.EvolutionResponse, error) {
	// Setup QueryForTemporalEvolution
	if req.Start == "" || req.End == "" {
		req.Start = "2020-02-19T04:30:00.000Z"
		req.End = "2020-02-19T06:30:00.000Z"
		req.Interval = "1m"
	}

	timeStamp := request.TimeStamp{"strict_date_optional_time", req.Start, req.End}
	range2 := request.Range{timeStamp}
	rangeObj := request.RangeObj{range2}
	filter := []request.RangeObj{rangeObj}

	b := request.Bool{filter}
	query := request.Query{b}

	// Setup AggregationsUtilization request
	histogram := request.DateHistogram{"@timestamp", req.Interval, "Y-M-d HH:mm", "+05:30"}
	sum := request.Sum{"flow.bytes"}
	myBandwidth := request.MyBandwidth{sum}
	aggsMyBandwidth := request.AggsMyBandwidth{myBandwidth}
	termsInclude := request.TermsInclude{"flow.service_name", ".*http.*"}
	byApps := request.ByApps{termsInclude, aggsMyBandwidth}

	termsExclude := request.TermsExclude{"flow.service_name", ".*http.*"}
	byOthers := request.ByOthers{termsExclude, aggsMyBandwidth}
	aggs := request.Aggs{byApps, byOthers}

	myDateHistogram := request.MyDateHistogram{histogram, aggs}
	aggregations := request.Aggregations{myDateHistogram}
	// Final request object
	temporalEvolutionRequest := request.TemporalEvolutionRequest{Size: 0, Aggregations: aggregations, Query: query}

	res, err := ts.Repository.QueryForTemporalEvolution(temporalEvolutionRequest)

	return parseResponse(res), err
}

func (ts TemporalEvolutionService) BuildTemporalSelectionData() (response.UiData, error) {
	// Setup QueryForTemporalEvolution
	timeStamp := request.TimeStamp{"strict_date_optional_time", "2020-02-19T04:30:00.000Z", "2020-02-19T06:30:00.000Z"}
	range2 := request.Range{timeStamp}
	rangeObj := request.RangeObj{range2}
	filter := []request.RangeObj{rangeObj}

	b := request.Bool{filter}
	query := request.Query{b}

	// Setup AggregationsUtilization request
	histogram := request.DateHistogram{"@timestamp", "1m", "Y-M-d HH:mm", "+05:30"}
	sum := request.Sum{"flow.bytes"}
	myBandwidth := request.MyBandwidth{sum}
	aggsMyBandwidth := request.AggsMyBandwidth{myBandwidth}

	myDateHistogram := request.MyDateHistogramSelection{histogram, aggsMyBandwidth}
	aggregations := request.AggregationsSelection{myDateHistogram}
	// Final request object
	temporalSelectionRequest := request.TemporalSelectionRequest{Size: 0, AggregationsSelection: aggregations, Query: query}

	res, err := ts.Repository.QueryForTemporalSelection(temporalSelectionRequest)
	return parseResponseSelection(res), err
}

func parseResponseSelection(evolutionResponse domain.TemporalSelectionResponse) response.UiData {
	labels := make([]string, 0)
	buckets := evolutionResponse.AggregationsSelection.MyDateHistogramSelection.Buckets
	bucketsData := make([]float32, 0)
	for _, bucket := range buckets {
		bucketsData = append(bucketsData, (bucket.MyBandWidth.Value*8/1000000)/60)
		labels = append(labels, bucket.KeyAsString)
	}
	evolutionObject := response.EvolutionObject{
		BorderWidth:     1,
		BorderColor:     "#738893",
		BackgroundColor: "#c5ccd3",
		Label:           "Bandwidth Consumption",
		Data:            bucketsData}
	return response.UiData{labels, []response.EvolutionObject{evolutionObject}}
}

func parseResponse(evolutionResponse domain.TemporalEvolutionResponse) response.EvolutionResponse {
	labels := make([]string, 0)
	buckets := evolutionResponse.Aggregations.MyDateHistogram.Buckets
	objectMap := make(map[string]response.EvolutionObject)
	for _, bucket := range buckets {
		labels = append(labels, bucket.KeyAsString)
		appBuckets := bucket.ByApps.Buckets
		for _, appBucket := range appBuckets {
			mObject := response.EvolutionObject{}
			if object, ok := objectMap[appBucket.Key]; ok {
				mObject = object
			} else {
				mObject = response.EvolutionObject{
					1,
					getRandomColor(),
					getRandomColor(),
					appBucket.Key,
					make([]float32, 0),
					0, 0}
			}
			mObject.Data = append(mObject.Data, (appBucket.MyBandWidth.Value*8/1000000)/60)
			mObject.BarData = mObject.BarData + (appBucket.MyBandWidth.Value*8/1000000)/60
			mObject.Count++
			objectMap[appBucket.Key] = mObject
		}

		otherBuckets := bucket.ByOthers.Buckets
		for _, otherBucket := range otherBuckets {
			mObject := response.EvolutionObject{}
			if object, ok := objectMap["others"]; ok {
				mObject = object
			} else {
				mObject = response.EvolutionObject{
					1,
					getRandomColor(),
					getRandomColor(),
					"other",
					make([]float32, 0),
					0, 0}
			}
			mObject.Data = append(mObject.Data, (otherBucket.MyBandWidth.Value*8/1000000)/60)
			mObject.BarData = mObject.BarData + (otherBucket.MyBandWidth.Value*8/1000000)/60
			mObject.Count++
			objectMap[otherBucket.Key] = mObject
		}
	}

	dataSets := make([]response.EvolutionObject, 0)
	barDataList := make([]float64, 0)
	barColor := make([]string, 0)
	barLabels := make([]string, 0)
	for _, value := range objectMap {
		dataSets = append(dataSets, value)
		barDataList = append(barDataList, math.Round((float64(value.BarData)/float64(value.Count))*100)/100)
		barColor = append(barColor, value.BackgroundColor)
		barLabels = append(barLabels, value.Label)
	}

	finalBarObject := response.FinalBarObject{16, "Application Bar", barDataList, barColor}
	barDataSets := make([]response.FinalBarObject, 0)
	barDataSets = append(barDataSets, finalBarObject)
	uiBarObject := response.UiBarObject{barLabels, barDataSets}
	return response.EvolutionResponse{response.UiData{Labels: labels, DataSets: dataSets}, uiBarObject}
}

func getRandomColor() string {
	chartColors := []string{
		"rgb(255, 99, 132)",
		"rgb(255, 159, 64)",
		"rgb(255, 205, 86)",
		"rgb(75, 192, 192)",
		"rgb(54, 162, 235)",
		"rgb(153, 102, 255)",
		"rgb(201, 203, 207)",
	}

	rand.Seed(time.Now().UnixNano())
	return chartColors[rand.Intn(7)]
}
