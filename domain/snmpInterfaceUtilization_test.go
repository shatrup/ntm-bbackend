package domain

import (
	"ntm-backend/dto/response"
	"reflect"
	"testing"
	"time"
)

func Test_interfaces_Load_Dto_data(t *testing.T) {
	//Arrange
	test := BucketsUtilization{
		Key:            "",
		DocCount:       0,
		OutOctets:      OutOctets{1024},
		TotalBandwidth: TotalBandwidth{2048},
		InOctets:       InOctets{512},
	}
	//Act
	expected := response.InterfacesLoadDetails{Transmit: 6.66, Receive: 3.333}
	actual := test.InterfacesLoadDto(60)

	//Assert
	if reflect.DeepEqual(actual,expected){
		t.Error(actual)
	}
}

func Test_interfaces_download_should_not_be_zero(t *testing.T) {
//Arrange
	test := MyGraph{Buckets: []Buckets{{
		AsDifferenceIn:  AsDifferenceIn{73880},
	},
	{
		AsDifferenceIn:  AsDifferenceIn{68395},
	},}}
	//Act
	expected := []float64{9.85,9.12}
	expected1 := 9.485
	actual, average := test.InterfacesDownloadDto(60)
	//Assert
	if reflect.DeepEqual(actual,expected) == false {
		t.Error(actual)
	}
	if reflect.DeepEqual(average,expected1) == false {
		t.Error(average)
	}
}

func Test_interfaces_upload_should_not_be_zero(t *testing.T) {
	//Arrange
	test := MyGraph{Buckets: []Buckets{{
		AsDifferenceOut:  AsDifferenceOut{73880},
	},
		{
			AsDifferenceOut:  AsDifferenceOut{68395},
		},}}
	//Act
	expected := []float64{9.85,9.12}
	expected1 := 9.485
	actual, average := test.InterfacesUploadDto(60)
	//Assert
	if reflect.DeepEqual(actual,expected) == false {
		t.Error(actual)
	}
	if reflect.DeepEqual(average,expected1) == false {
		t.Error(average)
	}
}

func Test_interfaces_total_bandwidth_should_not_be_zero(t *testing.T) {
	//Arrange
	test := MyGraph{Buckets: []Buckets{{
		AsDifferenceOut:  AsDifferenceOut{73880},
		AsDifferenceIn: AsDifferenceIn{73880},
	},
		{
			AsDifferenceOut:  AsDifferenceOut{68395},
			AsDifferenceIn: AsDifferenceIn{68395},
		},}}
	//Act
	expected := []float64{19.7,18.24}
	expected1 := 18.97
	actual, average := test.InterfacesTotalBandWidthDto( 60)
	//Assert
	if reflect.DeepEqual(actual, expected) == false {
		t.Error(actual)
	}
	if reflect.DeepEqual(average,expected1) == false {
		t.Error(actual)
	}
}

func Test_interfaces_date_should_not_be_nil(t *testing.T) {
	//Arrange
	parsedDate1, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:54:00")
	parsedDate2, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:55:00")
	test := MyGraph{Buckets: []Buckets{
		{
			KeyAsString: parsedDate1,
		},
		{
			KeyAsString: parsedDate2,
		},
	}}

	//Act
	actual := MyGraph.InterfacesDateDto(test)
	expectedTime1, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:54:00 +0000 UTC")
	expectedTime2, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:55:00 +0000 UTC")
	expected := []time.Time{expectedTime1, expectedTime2}

	//Assert
	if reflect.DeepEqual(actual, expected)  {
		t.Error(actual)
	}
}
