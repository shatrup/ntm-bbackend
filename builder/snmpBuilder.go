package builder

import (
	"github.com/sirupsen/logrus"
	"ntm-backend/domain"
)

type SnmpBuilder struct {
	Logger *logrus.Logger
}

func (b SnmpBuilder) BuildRange(startTime string, endTime string) map[string]interface{} {
	startDate := domain.Date{startTime}
	endDate := domain.Date{endTime}
	utcStartDate, _ := startDate.UTCDate()
	utcEndDate, _ := endDate.UTCDate()
	return map[string]interface{}{
		"@timestamp": map[string]interface{}{
			"format": "strict_date_optional_time",
			"gte":    utcStartDate,
			"lte":    utcEndDate,
		},
	}
}

func (b SnmpBuilder) BuildMatch(nasId string) map[string]interface{} {
	return map[string]interface{}{
		"system.sysName": nasId,
	}
}

func (b SnmpBuilder) BuildAggregationsForGraph(intervalTime string) map[string]interface{} {
	if intervalTime != "" {
		return map[string]interface{}{
			"get_nested_interfaces": map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "interfaces",
				},
				"aggs": BuildByInterfaceGraph(intervalTime),
			},
		}
	}
	return map[string]interface{}{}
}

func BuildByInterfaceGraph(intervalTime string) map[string]interface{} {
	return map[string]interface{}{
		"by_interfaces": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "interfaces.ifDescr.keyword",
			},
			"aggs": map[string]interface{}{
				"my_graph": map[string]interface{}{
					"date_histogram": BuildDateHistogram(intervalTime),
					"aggs":           BuildSerialDifference(),
				},
			},
		},
	}
}

func BuildSerialDifference() map[string]interface{} {
	return map[string]interface{}{
		"inOctets": map[string]interface{}{
			"avg": map[string]interface{}{
				"field": "interfaces.stats.ifHCInOctets",
			},
		},
		"outOctets": map[string]interface{}{
			"avg": map[string]interface{}{
				"field": "interfaces.stats.ifHCOutOctets",
			},
		},
		"as_difference_in": map[string]interface{}{
			"serial_diff": map[string]interface{}{
				"buckets_path": "inOctets",
				"lag":          1,
			},
		},
		"as_difference_out": map[string]interface{}{
			"serial_diff": map[string]interface{}{
				"buckets_path": "outOctets",
				"lag":          1,
			},
		},
	}
}

func BuildDateHistogram(intervalTime string) map[string]interface{} {
	return map[string]interface{}{
		"field":          "interfaces.ts",
		"fixed_interval": intervalTime,
	}
}

func (b SnmpBuilder) GetNestedStorage(intervalTime string) map[string]interface{} {
	return map[string]interface{}{
		"get_nested_storage": map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "storage",
			},
			"aggs": BuildDiskMemoryTotalAverage(),
		},
		"get_cpu_total_load_average": map[string]interface{}{
			"avg": map[string]interface{}{
				"field": "cpu.hrProcessorLoad",
			},
		},
		"graph_data_for_cpu": map[string]interface{}{
			"date_histogram": map[string]interface{}{
				"field":          "@timestamp",
				"fixed_interval": intervalTime,
				"format":         "Y-M-d HH:mm",
			},
			"aggs": map[string]interface{}{
				"get_cpu_average": map[string]interface{}{
					"avg": map[string]interface{}{
						"field": "cpu.hrProcessorLoad",
					},
				},
			},
		},
		"get_nested_storage_data": map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "storage",
			},
			"aggs": map[string]interface{}{
				"by_disk": map[string]interface{}{
					"terms": map[string]interface{}{
						"field":   "storage.hrStorageDescr.keyword",
						"include": "system disk",
					},
					"aggs": map[string]interface{}{
						"graph_data_for_disk_space": map[string]interface{}{
							"date_histogram": map[string]interface{}{
								"field":          "storage.ts",
								"fixed_interval": intervalTime,
							},
							"aggs": map[string]interface{}{
								"disk_average": map[string]interface{}{
									"avg": map[string]interface{}{
										"field": "storage.hrStorageUsed",
									},
								},
							},
						},
					},
				},
				"by_memory": map[string]interface{}{
					"terms": map[string]interface{}{
						"field":   "storage.hrStorageDescr.keyword",
						"include": "main memory",
					},
					"aggs": map[string]interface{}{
						"graph_data_for_memory": map[string]interface{}{
							"date_histogram": map[string]interface{}{
								"field":          "storage.ts",
								"fixed_interval": intervalTime,
							},
							"aggs": map[string]interface{}{
								"memory_average": map[string]interface{}{
									"avg": map[string]interface{}{
										"field": "storage.hrStorageUsed",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func BuildDiskMemoryTotalAverage() map[string]interface{} {
	return map[string]interface{}{
		"by_disk_storage": BuildDiskTotalAverage(),
		"By_memory_size":  BuildMemoryTotalAverage(),
	}
}

func BuildDiskTotalAverage() map[string]interface{} {
	return map[string]interface{}{
		"terms": map[string]interface{}{
			"field":   "storage.hrStorageDescr.keyword",
			"include": "system disk",
		},
		"aggs": map[string]interface{}{
			"Disk": map[string]interface{}{
				"avg": map[string]interface{}{
					"field": "storage.hrStorageUsed",
				},
			},
			"diskSize": map[string]interface{}{
				"avg": map[string]interface{}{
					"field": "storage.hrStorageSize",
				},
			},
		},
	}
}

func BuildMemoryTotalAverage() map[string]interface{} {
	return map[string]interface{}{
		"terms": map[string]interface{}{
			"field":   "storage.hrStorageDescr.keyword",
			"include": "main memory",
		},
		"aggs": map[string]interface{}{
			"Memory": map[string]interface{}{
				"avg": map[string]interface{}{
					"field": "storage.hrStorageUsed",
				},
			},
			"MemorySize": map[string]interface{}{
				"avg": map[string]interface{}{
					"field": "storage.hrStorageSize",
				},
			},
		},
	}
}

func (b SnmpBuilder) BuildInterfaceForUtilization() map[string]interface{} {
	return map[string]interface{}{
		"by_interfaces": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "interfaces.ifDescr.keyword",
			},
			"aggs": InterfaceAverage(),
		},
	}
}

func InterfaceAverage() map[string]interface{} {
	return map[string]interface{}{
		"inOctets": map[string]interface{}{
			"avg": map[string]interface{}{
				"field": "interfaces.stats.ifHCInOctets",
			},
		},
		"outOctets": map[string]interface{}{
			"avg": map[string]interface{}{
				"field": "interfaces.stats.ifHCOutOctets",
			},
		},
		"totalBandwidth": map[string]interface{}{
			"avg": map[string]interface{}{
				"field": "interfaces.ifSpeed",
			},
		},
	}
}
