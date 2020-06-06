package apps

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"ntm-backend/dto/request"
	"ntm-backend/dto/response"
	"ntm-backend/errs"
	"ntm-backend/service"
)

type ApiRouter struct {
	Logger     *logrus.Logger
	Service    service.TemporalEvolutionService
	snmService service.ISnmpService
}

func (api ApiRouter) RegisterApiRoutes(r *mux.Router, s string) {
	ur := r.PathPrefix(s).Subrouter()
	ur.HandleFunc("/temporalEvolution", api.TemporalEvolution).Methods(http.MethodPost)
	ur.HandleFunc("/temporalSelection", api.TemporalSelection).Methods(http.MethodPost)
	ur.HandleFunc("/snmpSystemSearch", api.SnmpSystemSearch).Methods(http.MethodPost)
	ur.HandleFunc("/snmpDeviceUtilization", api.SnmpDeviceUtilization).Methods(http.MethodPost)
	ur.HandleFunc("/snmpInterfacesUtilization", api.SnmpInterfacesUtilization).Methods(http.MethodPost)
	ur.HandleFunc("/snmpInterfaces", api.SnmpInterface).Methods(http.MethodPost)
}

func(api ApiRouter) RequestHandle(w http.ResponseWriter, r *http.Request) (request.SnmpRequest, bool) {
	var requestBody request.SnmpRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestBody)
	if err != nil {
		api.Logger.Println("Decoding error", err.Error())
		response.New(w, api.Logger).WithStatus(http.StatusBadRequest).WithResponse(errs.NewGeneralError(err))
	} else if bErrs := requestBody.Validate(); bErrs != nil {
			api.Logger.Println("Some error while valid request body Error: ", bErrs)
			response.New(w, api.Logger).WithStatus(http.StatusUnprocessableEntity).WithResponse(bErrs)
		} else {
			return requestBody, true
	}
	return request.SnmpRequest{}, false
}

func (api ApiRouter) SnmpInterface(w http.ResponseWriter, r *http.Request) {
	requestBody, isTrue  := api.RequestHandle(w,r)
	if isTrue {
		interfaceDetails, err := api.snmService.SnmpGraphDetails(requestBody.NasId, requestBody.StartDate, requestBody.EndDate)
		if err != nil {
			api.Logger.Println("Error Getting interface Details:", err.Error())
			response.New(w, api.Logger).WithStatus(http.StatusInternalServerError).WithResponse(err)
		} else {
			response.New(w, api.Logger).WithResponse(interfaceDetails)
		}
	}
}

func (api ApiRouter) SnmpInterfacesUtilization(w http.ResponseWriter, r *http.Request) {
	requestBody, isTrue  := api.RequestHandle(w,r)
	if isTrue {
		utilizationDetails, err := api.snmService.SnmpUtilizationDetails(requestBody.NasId, requestBody.StartDate, requestBody.EndDate)
		if err != nil {
			api.Logger.Println("Error while getting interface utilization details")
			response.New(w, api.Logger).WithStatus(http.StatusInternalServerError).WithResponse(err)
		} else {
			response.New(w, api.Logger).WithResponse(utilizationDetails)
		}
	}
}

func (api ApiRouter) SnmpDeviceUtilization(w http.ResponseWriter, r *http.Request) {
	requestBody, isTrue  := api.RequestHandle(w,r)
	if isTrue {
		deviceDetails, err := api.snmService.SnmpDeviceDetails(requestBody.NasId, requestBody.StartDate, requestBody.EndDate)
		if err != nil {
			api.Logger.Println("Error Getting Device Details")
			response.New(w, api.Logger).WithStatus(http.StatusInternalServerError).WithResponse(err)
		} else {
			response.New(w, api.Logger).WithResponse(deviceDetails)
		}
	}
}

func (api ApiRouter) SnmpSystemSearch(w http.ResponseWriter, r *http.Request) {
	requestBody, isTrue  := api.RequestHandle(w,r)
	if isTrue {
		systemDetails, err := api.snmService.SnmpSystemDetails(requestBody.NasId, requestBody.StartDate, requestBody.EndDate)
		if err != nil {
			api.Logger.Println("Error while getting system details")
			response.New(w, api.Logger).WithStatus(http.StatusInternalServerError).WithResponse(err)
		} else {
			response.New(w, api.Logger).WithResponse(systemDetails)
		}
	}
}

func (api ApiRouter) TemporalSelection(w http.ResponseWriter, r *http.Request) {
	response, err := api.Service.BuildTemporalSelectionData()
	if err != nil {
		api.Logger.Print("Error while returning response")
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

func (api ApiRouter) TemporalEvolution(w http.ResponseWriter, r *http.Request) {
	var temporalRequest request.TemporalRequest
	json.NewDecoder(r.Body).Decode(&temporalRequest)
	response, err := api.Service.BuildTemporalEvolutionData(temporalRequest)
	if err != nil {
		api.Logger.Print("Error while returning response")
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}
