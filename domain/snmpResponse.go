package domain

import (
	"fmt"
	"log"
	"math"
	"net"
	"ntm-backend/dto/response"
	"strconv"
	"time"
)

type AutoGenerated struct {
	ParentHits     ParentHits   `json:"hits"`
	TimedOut bool   `json:"timed_out"`
	Took     int    `json:"took"`
}

func (a AutoGenerated) GetSnmpSystemDetails() (response.DeviceDetails, error) {
	var deviceDetails response.DeviceDetails
	for _, hit := range a.ParentHits.Hits {
	deviceDetails.ModelNumber = hit.Source.System.SysDescr
	deviceDetails.TotalInterface = len(hit.Source.Interfaces)
	upTime, err := hit.Source.System.UpTime()
	if err != nil {
		return response.DeviceDetails{}, err
	}
	deviceDetails.CurrentUpTime = upTime
	memoryUsage, err := hit.Source.GetStorageUsage("main memory")
	if err != nil {
		return response.DeviceDetails{}, err
	}
	deviceDetails.MemoryUsage = memoryUsage
	diskUsage, err := hit.Source.GetStorageUsage("system disk")
	if err != nil {
		return response.DeviceDetails{}, err
	}
	deviceDetails.DiskUsage = diskUsage
	deviceDetails.InterfaceDetails = hit.Source.InterfacesDetailsDto()
	}
	var snmpResObject response.DeviceDetails
	snmpResObject = deviceDetails
	return snmpResObject, nil
}

type CPU struct {
	HrProcessorFrwID string `json:"hrProcessorFrwID"`
	HrProcessorLoad  string `json:"hrProcessorLoad"`
}
type Health struct {
	MtxrHITemperature string `json:"mtxrHITemperature"`
	MtxrHIVoltage     string `json:"mtxrHIVoltage"`
}
type Address struct {
	IPAdEntAddr    string `json:"ipAdEntAddr"`
	IPAdEntIfIndex string `json:"ipAdEntIfIndex"`
	IPAdEntNetMask string `json:"ipAdEntNetMask"`
}
type Stats struct {
	IfAlias              string `json:"ifAlias"`
	IfHCInBroadcastPkts  string `json:"ifHCInBroadcastPkts"`
	IfHCInMulticastPkts  string `json:"ifHCInMulticastPkts"`
	IfHCInOctets         string `json:"ifHCInOctets"`
	IfHCInUcastPkts      string `json:"ifHCInUcastPkts"`
	IfHCOutBroadcastPkts string `json:"ifHCOutBroadcastPkts"`
	IfHCOutMulticastPkts string `json:"ifHCOutMulticastPkts"`
	IfHCOutOctets        string `json:"ifHCOutOctets"`
	IfHCOutUcastPkts     string `json:"ifHCOutUcastPkts"`
	IfHighSpeed          string `json:"ifHighSpeed"`
	IfName               string `json:"ifName"`
}
type Interfaces struct {
	Address       Address   `json:"address,omitempty"`
	IfAdminStatus string    `json:"ifAdminStatus"`
	IfDescr       string    `json:"ifDescr"`
	IfInOctets    string    `json:"ifInOctets"`
	IfIndex       string    `json:"ifIndex"`
	IfLastChange  string    `json:"ifLastChange"`
	IfMtu         string    `json:"ifMtu"`
	IfOperStatus  string    `json:"ifOperStatus"`
	IfOutOctets   string    `json:"ifOutOctets"`
	IfPhysAddress string    `json:"ifPhysAddress"`
	IfSpeed       string    `json:"ifSpeed"`
	IfType        string    `json:"ifType"`
	Stats         Stats     `json:"stats"`
	Ts            time.Time `json:"ts"`
}

func (i Interfaces) GetAdminStatus() string {
	var adminStatus string
	if i.IfAdminStatus == "1" {
		adminStatus = "UP"
	} else if i.IfAdminStatus == "2" {
		adminStatus = "DOWN"
	} else {
		adminStatus = "TESTING"
	}
	return adminStatus
}
func (i Interfaces) GetOprStatus() string {
	var opStatus string
	if i.IfOperStatus == "1" {
		opStatus = "UP"
	} else if i.IfOperStatus == "2" {
		opStatus = "DOWN"
	} else {
		opStatus = "TESTING"
	}
	return opStatus
}
func (i Interfaces) GetPortType() string {
	var portType string
	if i.IfType == "1" {
		portType = "VPN"
	} else if i.IfType == "6" {
		portType = "Ethernet"
	}else if i.IfType == "71" {
		portType = "Wlan"
	}else if i.IfType == "135" {
		portType = "VLAN"
	}else if i.IfType == "209" {
		portType = "Bridge"
	} else {
		portType = "Unknown"
	}
	return portType
}

func (s Source) InterfacesDetailsDto() []response.InterfaceDetails {
	var interfacesDetailsDto []response.InterfaceDetails
	for _, ether := range s.Interfaces {
		if ether.Address.IPAdEntAddr == "" {
			continue
		} else if ether.IfDescr == "BhaiFi-VPN" {
			continue
		} else {
			sizeInByte, err := strconv.ParseFloat(ether.IfSpeed, 64)
			if err != nil {
				log.Fatal(err)
			}
			var sizeInBytes string
			sizeInBytes = bytesToAny(sizeInByte)
			addr := net.ParseIP(ether.Address.IPAdEntNetMask).To4()
			subnetInt , _ := net.IPv4Mask(addr[0], addr[1], addr[2], addr[3]).Size()
			subnet := strconv.Itoa(subnetInt)
			ethInfo := response.InterfaceDetails{AdminStatus:ether.GetAdminStatus(),OprationStatus:ether.GetOprStatus(),InterfaceName:ether.Stats.IfName,MacAddress:ether.IfPhysAddress,
				InterfaceType:ether.GetPortType(),NegotionSpeed:sizeInBytes,IpAddress: ether.Address.IPAdEntAddr,Subnet:subnet}
			interfacesDetailsDto = append(interfacesDetailsDto,ethInfo)
		}
	}
	return interfacesDetailsDto
}

type Storage struct {
	HrStorageAllocationFailures string `json:"hrStorageAllocationFailures"`
	HrStorageAllocationUnits    string `json:"hrStorageAllocationUnits"`
	HrStorageDescr              string `json:"hrStorageDescr"`
	HrStorageIndex              string `json:"hrStorageIndex"`
	HrStorageSize               string `json:"hrStorageSize"`
	HrStorageType               string `json:"hrStorageType"`
	HrStorageUsed               string `json:"hrStorageUsed"`
}

func (s Storage) DiskMemoryStorage (x string) (int, error)  {
		used, err := strconv.ParseInt(s.HrStorageUsed, 10, 64)
		if err != nil {
			return 0, err
		}
		size, err := strconv.ParseInt(s.HrStorageSize, 10, 64)
		if err != nil {
			return 0, err
		}
		y := (used * 100) / size
		if s.HrStorageDescr == x {
			return int(y), nil
		}
	return 0, err

}
type System struct {
	SysContact  string `json:"sysContact"`
	SysDescr    string `json:"sysDescr"`
	SysName     string `json:"sysName"`
	SysObjectID string `json:"sysObjectID"`
	SysUpTime   string `json:"sysUpTime"`
}
func (s System) UpTime () (string, error) {
	i, err := strconv.ParseInt(s.SysUpTime, 10, 64)
	if err != nil {
		return "", err
	}
	var output string
	uptime := i / 360000
	if uptime >= 1 {
		day := uptime / 24
		mod := uptime % 24;
		output = fmt.Sprintf("%d%s%d%s", day, " day ", mod, " Hours ")
	} else {
		output = fmt.Sprintf("%d%s", i, " Hours ")
	}
	return output, nil
}
type Source struct {
	Timestamp  time.Time    `json:"@timestamp"`
	CPU        []CPU        `json:"cpu"`
	Health     Health       `json:"health"`
	Interfaces []Interfaces `json:"interfaces"`
	SnmpHost   string       `json:"snmp.host"`
	Storage    []Storage    `json:"storage"`
	System     System       `json:"system"`
	Type       string       `json:"type"`
}

func (s Source) GetStorageUsage (storageName string ) (int, error) {
	if storageName != ""{
		for _, storage := range s.Storage {
			if storage.HrStorageDescr == storageName {
				storageUsage, err := storage.DiskMemoryStorage(storageName)
				if err != nil {
					return 0, err
				}
				return storageUsage, nil
			}
		}
	}
	return 0, nil
}
type Hits struct {
	ID     string      `json:"_id"`
	Index  string      `json:"_index"`
	Score  interface{} `json:"_score"`
	Source Source      `json:"_source"`
	Type   string      `json:"_type"`
	Sort   []int64     `json:"sort"`
}
type Total struct {
	Relation string `json:"relation"`
	Value    int    `json:"value"`
}
type ParentHits struct {
	Hits     []Hits      `json:"hits"`
	MaxScore interface{} `json:"max_score"`
	Total    Total       `json:"total"`
}


func Round(val float64, roundOn float64, places int ) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func bytesToAny(value float64) string {
	var suffixes [5]string
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"
	if value != 0 {
		base := math.Log(value)/math.Log(1000)
		getSize := Round(math.Pow(1000, base - math.Floor(base)), .5, 2)
		getSuffix := suffixes[int(math.Floor(base))]
		return (strconv.FormatFloat(getSize, 'f', -1, 64)+" "+string(getSuffix))
	}
	return ""

}