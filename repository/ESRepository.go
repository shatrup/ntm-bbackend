package repository

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"ntm-backend/domain"
	"ntm-backend/dto/request"
)

type ESRepository struct {
	Logger *logrus.Logger
}

func (es ESRepository) QueryForTemporalEvolution(request request.TemporalEvolutionRequest) (domain.TemporalEvolutionResponse, error) {
	//func (es ESRepository) QueryForTemporalEvolution(request request.TemporalEvolutionRequest) (interface{}, error) {

	response, err := es.esHttp(request, "/elastiflow-*/_search")
	if err != nil {
		es.Logger.Println("Error while making Request: ", err.Error())
		return domain.TemporalEvolutionResponse{}, err
	}
	//var temporalEvolutionResponse interface{}
	var temporalEvolutionResponse domain.TemporalEvolutionResponse
	err = json.NewDecoder(response.Body).Decode(&temporalEvolutionResponse)
	if err != nil {
		es.Logger.Println("Error while reading response: ", err.Error())
		return domain.TemporalEvolutionResponse{}, err
	}

	return temporalEvolutionResponse, nil
}

func (es ESRepository) QueryForTemporalSelection(request request.TemporalSelectionRequest) (domain.TemporalSelectionResponse, error) {
	response, err := es.esHttp(request, "/elastiflow-*/_search")
	if err != nil {
		es.Logger.Println("Error while making Request: ", err.Error())
		return domain.TemporalSelectionResponse{}, err
	}
	//var temporalEvolutionResponse interface{}
	var temporalSelectionResponse domain.TemporalSelectionResponse
	err = json.NewDecoder(response.Body).Decode(&temporalSelectionResponse)
	if err != nil {
		es.Logger.Println("Error while reading response: ", err.Error())
		return domain.TemporalSelectionResponse{}, err
	}

	return temporalSelectionResponse, nil
}

func (es ESRepository) esHttp(request interface{}, route string) (*http.Response, error) {
	baseUrl := "http://137.59.52.242:9200"
	endPoint := baseUrl + route
	jsonBytes, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodGet, endPoint, bytes.NewBuffer(jsonBytes))
	req.Header.Add("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}
