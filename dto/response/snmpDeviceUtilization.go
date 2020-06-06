package response

import "time"

type DeviceMemoryLoad struct {
	DeviceMemoryLoadDetails []DeviceMemoryLoadDetails `json:"deviceMemoryLoadDetails,omitempty"`
}

type DeviceMemoryLoadDetails struct {
	Date								[]time.Time 			`json:"key_as_string,omitempty"`
	Cpu									[]float64					`json:"cpu,omitempty"`
	Memory							[]float64					`json:"memory,omitempty"`
	Disk								[]float64					`json:"disk,omitempty"`
	AverageCpu					float64						`json:"averageCpu,omitempty"`
	AverageMemory 			float64						`json:"averageMemory,omitempty"`
	AverageDisk 				float64						`json:"averageDisk,omitempty"`
}