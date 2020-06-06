package domain

import (
	"fmt"
	"ntm-backend/dto/response"
	"reflect"
	"testing"
)

func Test_storage_usage_by_name_should_not_zero(t *testing.T) {
	test := Source{
		Storage:    []Storage{{
			HrStorageAllocationFailures: "",
			HrStorageAllocationUnits:    "",
			HrStorageDescr:              "main memory",
			HrStorageIndex:              "",
			HrStorageSize:               "131072",
			HrStorageType:               "",
			HrStorageUsed:               "27940",
		}, {
			HrStorageAllocationFailures: "",
			HrStorageAllocationUnits:    "",
			HrStorageDescr:              "system disk",
			HrStorageIndex:              "",
			HrStorageSize:               "131072",
			HrStorageType:               "",
			HrStorageUsed:               "20568",
		}},
	}
	value, err := test.GetStorageUsage("system disk")
	fmt.Println(value)
	if err != nil {
		t.Error(err)
	}
	if value == 0 {
		t.Error(value)
	}
}

func Test_up_time_should_not_be_nil(t *testing.T) {
	test := System{
		SysUpTime:   "1368474800",
	}
	value, err := test.UpTime()
	if err != nil {
		t.Error(err)
	}
	if value == "" {
		t.Error(value)
	}
}

func Test_interfaces_details_should_not_return_empty(t *testing.T) {
test := Source{
	Interfaces: []Interfaces{{
		Address:       Address{IPAdEntAddr: "10.4.0.49", IPAdEntNetMask: "255.255.255.240"},
		IfAdminStatus: "1",
		IfDescr:       "WAN1",
		IfOperStatus:  "2",
		IfPhysAddress: "6C:3B:6B:F5:A7:52",
		IfSpeed:       "100000000",
		IfType:        "6",
		Stats:         Stats{
			IfName:               "WAN1",
		},
	},
	{
		Address:       Address{IPAdEntAddr: "192.168.54.10",IPAdEntNetMask:"255.255.255.0"},
		IfAdminStatus: "1",
		IfDescr:       "WAN2",
		IfOperStatus:  "2",
		IfPhysAddress: "6C:3B:6B:F5:A7:53",
		IfSpeed:       "100000000",
		IfType:        "209",
		Stats:         Stats{
			IfName:               "WAN2",
		},
	}},
}

expected := []response.InterfaceDetails{{"UP", "DOWN"," WAN1", "6C:3B:6B:F5:A7:52", "Ethernet", "100 MB"," 10.4.0.49", "28"},{"UP", "DOWN", "WAN2", "6C:3B:6B:F5:A7:53", "Bridge", "100 MB", "192.168.54.10 ","24"}}
actual := test.InterfacesDetailsDto()
if reflect.DeepEqual(actual,expected) {
	t.Error(actual)
}
}