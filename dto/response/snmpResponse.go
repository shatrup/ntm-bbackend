package response

type DeviceDetails struct {
	TotalInterface		int 					`json:"totalInterface,omitempty"`
	MemoryUsage				int						`json:"memoryUsage,omitempty"`
	DiskUsage					int						`json:"diskUsage,omitempty"`
	ModelNumber      	string       	`json:"modelNumbers,omitempty"`
	SerialNumber			string				`json:"serialNumber,omitempty"`
	CurrentUpTime			string				`json:"currentUptime,omitempty"`
	InterfaceDetails  [] InterfaceDetails `json:"interfaceDetails,omitempty"`
}

type InterfaceDetails struct {
	AdminStatus				string 				`json:"adminStatus"`
	OprationStatus 		string 				`json:"oprationStatus"`
	InterfaceName			string 				`json:"interfaceName"`
	MacAddress				string 				`json:"macAddress"`
	InterfaceType 		string 				`json:"interfaceType"`
	NegotionSpeed			string 				`json:"negotionSpeed"`
	IpAddress					string 				`json:"ipAddress"`
	Subnet 						string 				`json:"subnet"`
}

