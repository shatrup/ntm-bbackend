package domain

import (
	"ntm-backend/dto/response"
	"time"
)

type SnmpDeviceUtilization struct {
	Aggregations Aggregations `json:"aggregations"`
}

func (s SnmpDeviceUtilization) GetSnmpDeviceDetails() (response.DeviceMemoryLoadDetails, error) {
	deviceUtilization := response.DeviceMemoryLoadDetails{}
	deviceUtilization.AverageCpu = s.Aggregations.GetCpuLoadAverage()
	deviceUtilization.AverageDisk = s.Aggregations.GetNestedStorage.ByDiskStorage.GetDiskUsagePercentage()
	deviceUtilization.AverageMemory = s.Aggregations.GetNestedStorage.ByMemorySize.GetMemoryUsagePercentage()
	deviceUtilization.Cpu = s.Aggregations.GraphDataForCPU.GetCpuGraphData()
	diskBucket := s.Aggregations.GetNestedStorage.ByDiskStorage.Buckets
	if (len(diskBucket)) != 0 {
		totalDiskSize := diskBucket[0].DiskSize.Value
		deviceUtilization.Disk = s.Aggregations.GetNestedStorageData.ByDisk.GetDiskGraphData(totalDiskSize)
	}
	memoryBucket := s.Aggregations.GetNestedStorage.ByMemorySize.Buckets
	if (len(memoryBucket)) != 0 {
		totalMemorySize := memoryBucket[0].MemorySize.Value
		deviceUtilization.Memory = s.Aggregations.GetNestedStorageData.ByMemory.GetMemoryGraphData(totalMemorySize)
	}
	deviceUtilization.Date = s.Aggregations.GetNestedStorageData.ByMemory.GetDateGraphData()
	return deviceUtilization, nil
}

type MemorySize struct {
	Value float64 `json:"value"`
}
type Memory struct {
	Value float64 `json:"value"`
}
type MemorySizeBuckets struct {
	Key        string     `json:"key"`
	DocCount   int        `json:"doc_count"`
	MemorySize MemorySize `json:"MemorySize"`
	Memory     Memory     `json:"memory"`
}

type ByMemorySize struct {
	DocCountErrorUpperBound int                 `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int                 `json:"sum_other_doc_count"`
	Buckets                 []MemorySizeBuckets `json:"buckets"`
}


func (m ByMemorySize) GetMemoryUsagePercentage() float64 {
	if (len(m.Buckets)) != 0 && m.Buckets[0].MemorySize.Value != 0 {
		mAverage := Round(m.Buckets[0].Memory.Value*100/m.Buckets[0].MemorySize.Value, .5, 2)
		return mAverage
	}
	return 0
}

type DiskSize struct {
	Value float64 `json:"value"`
}
type Disk struct {
	Value float64 `json:"value"`
}
type DiskSizeBuckets struct {
	Key      string   `json:"key"`
	DocCount int      `json:"doc_count"`
	DiskSize DiskSize `json:"diskSize"`
	Disk     Disk     `json:"disk"`
}

type ByDiskStorage struct {
	DocCountErrorUpperBound int               `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int               `json:"sum_other_doc_count"`
	Buckets                 []DiskSizeBuckets `json:"buckets"`
}


func (d ByDiskStorage) GetDiskUsagePercentage() float64 {
	if (len(d.Buckets)) != 0 && d.Buckets[0].DiskSize.Value != 0 {
		mAverage := Round(d.Buckets[0].Disk.Value*100/d.Buckets[0].DiskSize.Value, .5, 2)
		return mAverage
	}
	return 0
}

type GetNestedStorage struct {
	DocCount      int           `json:"doc_count"`
	ByMemorySize  ByMemorySize  `json:"by_memory_size"`
	ByDiskStorage ByDiskStorage `json:"by_disk_storage"`
}
type GetCPUAverage struct {
	Value float64 `json:"value"`
}
type CpuBuckets struct {
	KeyAsString   string        `json:"key_as_string"`
	Key           int64         `json:"key"`
	DocCount      int           `json:"doc_count"`
	GetCPUAverage GetCPUAverage `json:"get_cpu_average"`
}

type GraphDataForCPU struct {
	Buckets []CpuBuckets `json:"buckets"`
}

func (g GraphDataForCPU) GetCpuGraphData() []float64 {
	var cpuGraphData []float64
	for _, inter := range g.Buckets {
		receive := inter.GetCPUAverage.Value
		cpuGraphData = append(cpuGraphData, receive)
	}
	return cpuGraphData
}

type DiskAverage struct {
	Value float64 `json:"value"`
}
type DiskSpaceBuckets struct {
	KeyAsString time.Time   `json:"key_as_string"`
	Key         int64       `json:"key"`
	DocCount    int         `json:"doc_count"`
	DiskAverage DiskAverage `json:"disk_average"`
}
type GraphDataForDiskSpace struct {
	Buckets []DiskSpaceBuckets `json:"buckets"`
}

type DiskBuckets struct {
	Key                   string                `json:"key"`
	DocCount              int                   `json:"doc_count"`
	GraphDataForDiskSpace GraphDataForDiskSpace `json:"graph_data_for_disk_space"`
}
type ByDisk struct {
	DocCountErrorUpperBound int           `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int           `json:"sum_other_doc_count"`
	Buckets                 []DiskBuckets `json:"buckets"`
}

func (g ByDisk) GetDiskGraphData(TotalDiskSize float64) []float64 {
	var diskGraphData []float64
	if (len(g.Buckets)) != 0 {
		for _, inter := range g.Buckets {
			if (len(inter.GraphDataForDiskSpace.Buckets)) != 0 {
				for _, space := range inter.GraphDataForDiskSpace.Buckets {
					receive := Round(space.DiskAverage.Value*100/TotalDiskSize, .5, 2)
					diskGraphData = append(diskGraphData, receive)
				}
			}
		}
		return diskGraphData
	}
	return nil
}

type MemoryAverage struct {
	Value float64 `json:"value"`
}
type ForMemoryBuckets struct {
	KeyAsString   time.Time     `json:"key_as_string"`
	Key           int64         `json:"key"`
	DocCount      int           `json:"doc_count"`
	MemoryAverage MemoryAverage `json:"memory_average"`
}
type GraphDataForMemory struct {
	Buckets []ForMemoryBuckets `json:"buckets"`
}

type MemoryBuckets struct {
	Key                string             `json:"key"`
	DocCount           int                `json:"doc_count"`
	GraphDataForMemory GraphDataForMemory `json:"graph_data_for_memory"`
}
type ByMemory struct {
	DocCountErrorUpperBound int             `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int             `json:"sum_other_doc_count"`
	Buckets                 []MemoryBuckets `json:"buckets"`
}

func (g ByMemory) GetMemoryGraphData(TotalMemorySize float64) []float64 {
	var memoryGraphData []float64
	if (len(g.Buckets)) != 0 && TotalMemorySize != 0 {
		for _, inter := range g.Buckets {
			if (len(inter.GraphDataForMemory.Buckets)) != 0 {
				for _, mem := range inter.GraphDataForMemory.Buckets {
					receive := Round(mem.MemoryAverage.Value*100/TotalMemorySize, .5, 2)
					memoryGraphData = append(memoryGraphData, receive)
				}
			}
		}
		return memoryGraphData
	}
return nil
}

func (g ByMemory) GetDateGraphData() []time.Time {
	var DateGraphData []time.Time
	if (len(g.Buckets)) != 0 {
		for _, inter := range g.Buckets {
			if (len(inter.GraphDataForMemory.Buckets)) != 0 {
				for _, dates := range inter.GraphDataForMemory.Buckets {
					receive := dates.KeyAsString.Add(time.Hour * 5 + time.Minute *30)
					DateGraphData = append(DateGraphData, receive)
				}
			}
		}
		return DateGraphData
	}
	return nil
}

type GetNestedStorageData struct {
	DocCount int      `json:"doc_count"`
	ByDisk   ByDisk   `json:"by_disk"`
	ByMemory ByMemory `json:"by_memory"`
}
type GetCPUTotalLoadAverage struct {
	Value float64 `json:"value"`
}
type Aggregations struct {
	GetNestedStorage       GetNestedStorage       `json:"get_nested_storage"`
	GraphDataForCPU        GraphDataForCPU        `json:"graph_data_for_cpu"`
	GetNestedStorageData   GetNestedStorageData   `json:"get_nested_storage_data"`
	GetCPUTotalLoadAverage GetCPUTotalLoadAverage `json:"get_cpu_total_load_average"`
}

func (a Aggregations) GetCpuLoadAverage() float64 {
	cAverage := Round(a.GetCPUTotalLoadAverage.Value, .5, 2)
	return cAverage
}
