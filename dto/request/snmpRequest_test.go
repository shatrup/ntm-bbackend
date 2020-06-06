package request

import "testing"

func Test_SnmpRequest(t *testing.T)  {
	data := SnmpRequest{
		NasId: "d5d2",
		StartDate: "",
		EndDate: "",
	}
	err :=SnmpRequest.Validate(data)
	if err == nil {
		t.Error(err)
	}
}

func Test_SnmpRequest2(t *testing.T)  {
	data := SnmpRequest{
		NasId: "test_39_1",
		StartDate: "0",
		EndDate: "2020-05-12 12:00:01",
	}
	err :=SnmpRequest.Validate(data)
	if err == nil {
		t.Error(err)
	}
}

func Test_SnmpRequest3(t *testing.T)  {
	data := SnmpRequest{
		NasId: "test_39_1",
		StartDate: "2020-05-13 12:00:00",
		EndDate: "2020-05-12 12:00:01",
	}
	err :=SnmpRequest.Validate(data)
	if err == nil {
		t.Error(err)
	}
}