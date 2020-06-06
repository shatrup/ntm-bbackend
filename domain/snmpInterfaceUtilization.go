package domain

import (
	"ntm-backend/dto/response"
	"time"
)

type SnmpInterfaceUtilization struct {
	Aggregations AggregationsUtilization `json:"aggregations"`
}

func (s SnmpInterfaceUtilization) GetSnmpUtilizationDetails(second int) ([]response.InterfacesLoadDetails, error)  {
	interfaceLoad := make([]response.InterfacesLoadDetails, 0)
	if (len(s.Aggregations.GetNestedInterfaces.ByInterfaces.BucketsUtilization)) != 0 {
		for _, value := range s.Aggregations.GetNestedInterfaces.ByInterfaces.BucketsUtilization {
			if value.Key == "BhaiFi-VPN" {
				continue
			} else {
				totalLoad := value.InterfacesLoadDto(second)
				interfaceLoad = append(interfaceLoad, totalLoad)
			}
		}
	}
	return interfaceLoad, nil
}

type SnmpInterfaceGraph struct {
	AggregationsGraph AggregationsGraph `json:"aggregations"`
}

func (s SnmpInterfaceGraph) GetSnmpGraphDetails(second float64) ([]response.InterfaceGraphDetails, error) {
	interfaceGraphDetails := make([]response.InterfaceGraphDetails, 0)
	var interfaceGraph response.InterfaceGraphDetails
	if len(s.AggregationsGraph.GetNestedInterfaces.ByInterfaces.GraphBuckets) != 0 {
		for _, value := range s.AggregationsGraph.GetNestedInterfaces.ByInterfaces.GraphBuckets {
			interfaceGraph.InterfaceName = value.Key
			download, dAverage := value.MyGraph.InterfacesDownloadDto(second)
			interfaceGraph.Download = download
			interfaceGraph.AverageReceive = dAverage
			upload, uAverage := value.MyGraph.InterfacesUploadDto(second)
			interfaceGraph.Upload = upload
			interfaceGraph.AverageTransmit = uAverage
			totalBandWidth, tAverage := value.MyGraph.InterfacesTotalBandWidthDto(second)
			interfaceGraph.TotalBandWidth = totalBandWidth
			interfaceGraph.AverageBandwidth = tAverage
			interfaceGraph.Date = value.MyGraph.InterfacesDateDto()
			interfaceGraphDetails = append(interfaceGraphDetails, interfaceGraph)
		}
	}
	return interfaceGraphDetails, nil
}

type OutOctets struct {
	Value float64 `json:"value"`
}
type TotalBandwidth struct {
	Value float64 `json:"value"`
}
type InOctets struct {
	Value float64 `json:"value"`
}
type AsDifferenceIn struct {
	Value float64 `json:"value"`
}
type AsDifferenceOut struct {
	Value float64 `json:"value"`
}
type Buckets struct {
	KeyAsString     time.Time       `json:"key_as_string"`
	Key             int64           `json:"key"`
	DocCount        int             `json:"doc_count"`
	OutOctets       OutOctets       `json:"outOctets"`
	TotalBandwidth  TotalBandwidth  `json:"totalBandwidth"`
	InOctets        InOctets        `json:"inOctets"`
	AsDifferenceIn  AsDifferenceIn  `json:"as_difference_in,omitempty"`
	AsDifferenceOut AsDifferenceOut `json:"as_difference_out,omitempty"`
}

type BucketsUtilization struct {
	Key            string         `json:"key"`
	DocCount       int            `json:"doc_count"`
	OutOctets      OutOctets      `json:"outOctets"`
	TotalBandwidth TotalBandwidth `json:"totalBandwidth"`
	InOctets       InOctets       `json:"inOctets"`
}

func (b BucketsUtilization) InterfacesLoadDto(second int) response.InterfacesLoadDetails {
	var receive float64
	var transmit float64
	if b.TotalBandwidth.Value != 0.0 {
		receive = Round((b.InOctets.Value*8*100)/(float64(second)*b.TotalBandwidth.Value), .5, 4)
		transmit = Round((b.OutOctets.Value*8*100)/(float64(second)*b.TotalBandwidth.Value), .5, 4)
	} else {
		receive = 0.0
		transmit = 0.0
	}
	return response.InterfacesLoadDetails{InterfaceName: b.Key, Transmit: transmit, Receive: receive}
}

type MyGraph struct {
	Buckets []Buckets `json:"buckets"`
}

func (m MyGraph) InterfacesDownloadDto(seconds float64) ([]float64, float64) {
	var interfacesDownload []float64
	var averageSumReceive float64
	var average float64
	var bytes = 1000
	if len(m.Buckets) > 0 {
		for _, inter := range m.Buckets {
			receive := Round(((inter.AsDifferenceIn.Value*8)/float64(bytes))/(seconds), .5, 2)
			averageSumReceive = averageSumReceive + receive
			interfacesDownload = append(interfacesDownload, receive)
		}
		average = Round(averageSumReceive/float64(len(m.Buckets)), .5, 3)
	}
	averageSumReceive = 0
	return interfacesDownload, average
}
func (m MyGraph) InterfacesUploadDto(seconds float64) ([]float64, float64) {
	var interfacesUpload []float64
	var averageSumTransmit float64
	var average float64
	var bytes = 1000
	if len(m.Buckets) > 0 {
		for _, inter := range m.Buckets {
			receive := Round(((inter.AsDifferenceOut.Value*8)/float64(bytes))/(seconds), .5, 2)
			averageSumTransmit = averageSumTransmit + receive
			interfacesUpload = append(interfacesUpload, receive)
		}
		average = Round(averageSumTransmit/float64(len(m.Buckets)), .5, 3)
	}
	averageSumTransmit = 0
	return interfacesUpload, average
}
func (m MyGraph) InterfacesTotalBandWidthDto(seconds float64) ([]float64, float64) {
	var interfacesTotalBandWidth []float64
	var averageSumBandwidth float64
	var average float64
	var bytes = 1000
	if len(m.Buckets) > 0 {
		for _, inter := range m.Buckets {
			receive := Round((((inter.AsDifferenceIn.Value*8)/float64(bytes))/(seconds))+(((inter.AsDifferenceOut.Value*8)/float64(bytes))/(seconds)), .5, 2)
			averageSumBandwidth = averageSumBandwidth + receive
			interfacesTotalBandWidth = append(interfacesTotalBandWidth, receive)
		}
		average = Round(averageSumBandwidth/float64(len(m.Buckets)), .5, 3)
	}
	averageSumBandwidth = 0
	return interfacesTotalBandWidth, average
}
func (m MyGraph) InterfacesDateDto() []time.Time {
	var interfacesDate []time.Time
	for _, inter := range m.Buckets {
		interfacesDate = append(interfacesDate, inter.KeyAsString.Add(time.Hour*5+time.Minute*30))
	}
	return interfacesDate
}

type GraphBuckets struct {
	Key      string  `json:"key"`
	DocCount int     `json:"doc_count"`
	MyGraph  MyGraph `json:"my_graph"`
}
type ByInterfacesUtilization struct {
	DocCountErrorUpperBound int                  `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int                  `json:"sum_other_doc_count"`
	BucketsUtilization      []BucketsUtilization `json:"buckets"`
}

type GetNestedInterfacesUtilization struct {
	DocCount     int                     `json:"doc_count"`
	ByInterfaces ByInterfacesUtilization `json:"by_interfaces"`
}

type AggregationsUtilization struct {
	GetNestedInterfaces GetNestedInterfacesUtilization `json:"get_nested_interfaces"`
}

type ByInterfacesGraph struct {
	DocCountErrorUpperBound int            `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int            `json:"sum_other_doc_count"`
	GraphBuckets            []GraphBuckets `json:"buckets"`
}

type GetNestedInterfacesGraph struct {
	DocCount     int               `json:"doc_count"`
	ByInterfaces ByInterfacesGraph `json:"by_interfaces"`
}
type AggregationsGraph struct {
	GetNestedInterfaces GetNestedInterfacesGraph `json:"get_nested_interfaces"`
}
