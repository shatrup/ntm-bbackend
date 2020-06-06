package apps

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"ntm-backend/dto/response"
	service2 "ntm-backend/mocks/service"
	"ntm-backend/service"
	"strings"
	"testing"
	"time"
)

var router *mux.Router

var apiRouterService *service2.MockISnmpService

func setupSubTest(t *testing.T) func(t *testing.T) {
	t.Log("setup sub test")
	mockCtrl := gomock.NewController(t)

	apiRouterService = service2.NewMockISnmpService(mockCtrl)
	router = mux.NewRouter()
	apiRouter := ApiRouter{
		Logger:     logrus.New(),
		Service:    service.TemporalEvolutionService{},
		snmService: apiRouterService,
	}
	// server side
	apiRouter.RegisterApiRoutes(router, "/")

	return func(t *testing.T) {
		router = nil
		apiRouterService = nil
		defer mockCtrl.Finish()
	}

}

func Test_api_router_should_return_with_400(t *testing.T) {
	// Arrange
	tearDownSubTest := setupSubTest(t)
	defer tearDownSubTest(t)
	request := `{
	"asfjdas",
asdaslkf,
	}`
	testUrl := []string{"/snmpInterfaces", "/snmpInterfacesUtilization", "/snmpDeviceUtilization", "/snmpSystemSearch"}
	for _, unitUrl := range testUrl {
		// client side
		req, _ := http.NewRequest(http.MethodPost, unitUrl, strings.NewReader(request))
		response := httptest.NewRecorder()

		// Act
		router.ServeHTTP(response, req)

		// Assert
		expected := 400

		if expected != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", expected, response.Code)
		}
	}
}

func Test_api_router_should_return_with_422_invalid_start_date(t *testing.T) {
	// Arrange
	tearDownSubTest := setupSubTest(t)
	defer tearDownSubTest(t)
	// client side
	request := `{
		"nas_id": "of_111_1",
		"start_date": "",
		"end_date": ""
	}`
	testUrl := []string{"/snmpInterfaces", "/snmpInterfacesUtilization", "/snmpDeviceUtilization", "/snmpSystemSearch"}
	for _, unitUrl := range testUrl {
		req, _ := http.NewRequest(http.MethodPost, unitUrl, strings.NewReader(request))
		response := httptest.NewRecorder()

		// Act
		router.ServeHTTP(response, req)

		// Assert
		expected := 422

		fmt.Println("BODY::: ", response.Body)
		if expected != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", expected, response.Code)
		}
	}

}

func Test_api_router_should_return_with_200_and_no_data(t *testing.T) {
	// Arrange
	tearDownSubTest := setupSubTest(t)
	defer tearDownSubTest(t)

	request := `{
		"nas_id": "of_111_1",
		"start_date": "2020-05-19 13:00:00",
		"end_date": "2020-05-19 13:05:00"
	}`

	testSetup := []struct {
		url              string
		expectedResponse string
		responseCode     int
	}{
		{
			"/snmpInterfaces",
			"[]",
			200,
		},
		{
			"/snmpInterfacesUtilization",
			"[]",
			200,
		},
		{
			"/snmpDeviceUtilization",
			"{}",
			200,
		},
		{"/snmpSystemSearch",
			"{}",
			200,
		},
	}

	apiRouterService.EXPECT().SnmpGraphDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return([]response.InterfaceGraphDetails{}, nil).Times(1)
	apiRouterService.EXPECT().SnmpUtilizationDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return([]response.InterfacesLoadDetails{}, nil)
	apiRouterService.EXPECT().SnmpDeviceDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(response.DeviceMemoryLoadDetails{}, nil)
	apiRouterService.EXPECT().SnmpSystemDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(response.DeviceDetails{}, nil)

	for _, test := range testSetup {
		req, _ := http.NewRequest(http.MethodPost, test.url, strings.NewReader(request))
		res := httptest.NewRecorder()

		// Act
		router.ServeHTTP(res, req)

		// Assert
		s := strings.TrimRight(res.Body.String(), "\n")

		if s != test.expectedResponse {
			t.Errorf("Expected response %v , Got %v\n", test.expectedResponse, s)
		}

		if test.responseCode != res.Code {
			t.Errorf("Expected response code %d. Got %d\n", test.responseCode, res.Code)
		}
	}
}

func Test_api_router_should_return_with_200_and_data(t *testing.T) {
	// Arrange
	tearDownSubTest := setupSubTest(t)
	defer tearDownSubTest(t)
	// client side
	request := `{
		"nas_id": "of_111_1",
		"start_date": "2020-05-19 13:00:00",
		"end_date": "2020-05-19 13:05:00"
	}`
	SnmpSystemDetailsResponse := `{"totalInterface":16,"memoryUsage":56,"diskUsage":14,"modelNumbers":"IN_RouterOd-5PE","serialNumber":"25631_745","currentUptime":"2 days 23 hours","interfaceDetails":[{"adminStatus":"1","oprationStatus":"6","interfaceName":"Wan1","macAddress":"86:k6:6O:4G","interfaceType":"Ethernet","negotionSpeed":"10000","ipAddress":"192.56.89.12","subnet":"23"}]}`
	snmpGraphDetailsResponse := `[{"interfaceName":"InterfaceName","key_as_string":null,"upload":null,"download":null,"totalBandWidth":null,"averageBandwidth":0,"averageReceive":0,"averageTransmit":0}]`
	SnmpDeviceDetailsResponse :=`{"cpu":[91.63,74.63],"memory":[91.63,74.63],"disk":[91.63,74.63],"averageCpu":56.52,"averageMemory":41.63,"averageDisk":63.45}`
	SnmpUtilizationDetailsResponse := `[{"interfaceName":"Wan1","transmit":56,"receive":84},{"interfaceName":"wan2","transmit":63.5,"receive":89.6}]`

	testSetup := []struct {
		url              string
		expectedResponse string
		responseCode     int
	}{
		{
			"/snmpInterfaces",
			snmpGraphDetailsResponse,
			200,
		},
		{
			"/snmpInterfacesUtilization",
			SnmpUtilizationDetailsResponse,
			200,
		},
		{
			"/snmpDeviceUtilization",
			SnmpDeviceDetailsResponse,
			200,
		},
		{"/snmpSystemSearch",
			SnmpSystemDetailsResponse,
			200,
		},
	}

	SnmpSystemDetailsResponse1 := response.DeviceDetails{
		TotalInterface: 16,
		MemoryUsage:    56,
		DiskUsage:      14,
		ModelNumber:    "IN_RouterOd-5PE",
		SerialNumber:   "25631_745",
		CurrentUpTime:  "2 days 23 hours",
		InterfaceDetails: []response.InterfaceDetails{{
			AdminStatus:    "1",
			OprationStatus: "6",
			InterfaceName:  "Wan1",
			MacAddress:     "86:k6:6O:4G",
			InterfaceType:  "Ethernet",
			NegotionSpeed:  "10000",
			IpAddress:      "192.56.89.12",
			Subnet:         "23",
		}},
	}
	snmpGraphDetailsResponse1 := []response.InterfaceGraphDetails{{
		InterfaceName:    "InterfaceName",
		Date:             nil,
		Upload:           nil,
		Download:         nil,
		TotalBandWidth:   nil,
		AverageBandwidth: 0,
		AverageReceive:   0,
		AverageTransmit:  0,
	}}
	SnmpDeviceDetailsResponse1 := response.DeviceMemoryLoadDetails{
		Date:          []time.Time{},
		Cpu:           []float64{91.63, 74.63},
		Memory:        []float64{91.63, 74.63},
		Disk:          []float64{91.63, 74.63},
		AverageCpu:    56.52,
		AverageMemory: 41.63,
		AverageDisk:   63.45,
	}
	SnmpUtilizationDetailsResponse1 := []response.InterfacesLoadDetails{{
		InterfaceName: "Wan1",
		Transmit:      56.0,
		Receive:       84.0,
	},{
		InterfaceName: "wan2",
		Transmit:      63.5,
		Receive:       89.6,
	}}

	apiRouterService.EXPECT().SnmpGraphDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(snmpGraphDetailsResponse1, nil)
	apiRouterService.EXPECT().SnmpUtilizationDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(SnmpUtilizationDetailsResponse1, nil)
	apiRouterService.EXPECT().SnmpDeviceDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(SnmpDeviceDetailsResponse1, nil)
	apiRouterService.EXPECT().SnmpSystemDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(SnmpSystemDetailsResponse1, nil)
	//testUrl := []string{"/snmpInterfaces", "/snmpInterfacesUtilization", "/snmpDeviceUtilization", "/snmpSystemSearch"}
	for _, unitUrl := range testSetup {
		req, _ := http.NewRequest(http.MethodPost, unitUrl.url, strings.NewReader(request))
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)
		// Assert
		s := strings.TrimRight(resp.Body.String(), "\n")
		if s != unitUrl.expectedResponse {
			t.Errorf("Expected response %v , Got %v\n", unitUrl.expectedResponse, s)
		}
		if unitUrl.responseCode != resp.Code {
			t.Errorf("Expected response code %d. Got %d\n", unitUrl.responseCode , resp.Code)
		}
	}
}

func Test_api_router_should_return_with_500(t *testing.T) {
	// Arrange
	tearDownSubTest := setupSubTest(t)
	defer tearDownSubTest(t)
	// client side
	request := `{
		"nas_id": "of_111_1",
		"start_date": "2020-05-19 13:00:00",
		"end_date": "2020-05-19 13:05:00"
	}`
	apiRouterService.EXPECT().SnmpGraphDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(nil, errors.New("internal error "))
	apiRouterService.EXPECT().SnmpUtilizationDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(nil, errors.New("internal error "))
	apiRouterService.EXPECT().SnmpDeviceDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(response.DeviceMemoryLoadDetails{}, errors.New("internal error "))
	apiRouterService.EXPECT().SnmpSystemDetails("of_111_1", "2020-05-19 13:00:00", "2020-05-19 13:05:00").Return(response.DeviceDetails{}, errors.New("internal error "))

	testUrl := []string{"/snmpInterfaces", "/snmpInterfacesUtilization", "/snmpDeviceUtilization", "/snmpSystemSearch"}
	for _, unitUrl := range testUrl {
		req, _ := http.NewRequest(http.MethodPost, unitUrl, strings.NewReader(request))
		response := httptest.NewRecorder()

		// Act
		router.ServeHTTP(response, req)

		// Assert
		expected := 500

		if expected != response.Code {
			t.Errorf("Expected response code %d. Got %d\n", expected, response.Code)
		}
	}
}
