package response

import "time"

type InterfacesLoad struct {
	InterfacesLoadDetails []InterfacesLoadDetails `json:"interfacesLoadDetails`
	InterfaceGraphDetails	[]InterfaceGraphDetails	`json:"interfaceGraphDetails"`
}

type InterfacesLoadDetails struct {
	InterfaceName			string 				`json:"interfaceName"`
	Transmit					float64				`json:"transmit"`
	Receive						float64				`json:"receive"`
}

type InterfaceGraphDetails struct {
	InterfaceName				string 						`json:"interfaceName"`
	Date								[]time.Time 			`json:"key_as_string"`
	Upload							[]float64					`json:"upload"`
	Download						[]float64					`json:"download"`
	TotalBandWidth			[]float64					`json:"totalBandWidth"`
	AverageBandwidth		float64						`json:"averageBandwidth"`
	AverageReceive 			float64						`json:"averageReceive"`
	AverageTransmit 		float64						`json:"averageTransmit"`
}

