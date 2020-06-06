package domain

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_DiskSize_should_not_be_zero(t *testing.T) {
	test := ByDiskStorage{Buckets: []DiskSizeBuckets{{
		Key:      "",
		DocCount: 0,
		DiskSize: DiskSize{Value: 1024},
		Disk:     Disk{Value: 750},
	}}}
	value := ByDiskStorage.GetDiskUsagePercentage(test)
	if value == 0 {
		t.Error(value)
	}
}

func Test_DiskSize_Bucket_should_not_be_nil(t *testing.T) {
	test := ByDiskStorage{Buckets: []DiskSizeBuckets{{}}}
	value := ByDiskStorage.GetDiskUsagePercentage(test)
	if value != 0 {
		t.Error(value)
	}
}

func Test_MemorySize_should_not_be_zero(t *testing.T) {
	test := ByMemorySize{Buckets: []MemorySizeBuckets{{
		Key:        "",
		DocCount:   0,
		MemorySize: MemorySize{1024},
		Memory:     Memory{512},
	}}}
	value := ByMemorySize.GetMemoryUsagePercentage(test)
	if value == 0 {
		t.Error(value)
	}
}

func Test_Memory_Bucket_should_not_be_nil(t *testing.T) {
	test := ByMemorySize{Buckets: []MemorySizeBuckets{{
		Key:        "",
		DocCount:   0,
		MemorySize: MemorySize{0},
		Memory:     Memory{0},
	}}}
	value := ByMemorySize.GetMemoryUsagePercentage(test)
	if value != 0 {
		t.Error(value)
	}
}

func Test_GetCpuGraphData_should_return_slice_of_cpu_averages(t *testing.T) {
	test := GraphDataForCPU{Buckets: []CpuBuckets{{
		GetCPUAverage: GetCPUAverage{
			Value: 48.5,
		},
	}, {
		GetCPUAverage: GetCPUAverage{
			Value: 28.0,
		},
	}}}
	expected := []float64{48.5, 28.0}
	actual := GraphDataForCPU.GetCpuGraphData(test)
	if reflect.DeepEqual(actual, expected) == false {
		t.Error("Not matching expected")
	}
}

func Test_cpu_Bucket_should_not_be_nil(t *testing.T) {
	test := GraphDataForCPU{Buckets: nil}
	actual := GraphDataForCPU.GetCpuGraphData(test)
	var expected []float64
	if reflect.DeepEqual(actual, expected) == false {
		t.Error("Should return empty slice")
	}
}

func Test_Date_Bucket_should_not_be_nil(t *testing.T) {
	// AAA test structure
	// Arrange
	parsedDate1, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:54:00")
	parsedDate2, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:55:00")
	test := ByMemory{
		DocCountErrorUpperBound: 0,
		SumOtherDocCount:        0,
		Buckets: []MemoryBuckets{{
			Key:      "",
			DocCount: 0,
			GraphDataForMemory: GraphDataForMemory{Buckets: []ForMemoryBuckets{{
				KeyAsString:   parsedDate1,
				Key:           0,
				DocCount:      0,
				MemoryAverage: MemoryAverage{},
			}, {
				KeyAsString:   parsedDate2,
				Key:           0,
				DocCount:      0,
				MemoryAverage: MemoryAverage{},
			}}},
		}},
	}

	// Act
	actual := ByMemory.GetDateGraphData(test)

	expectedTime1, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:54:00 +0000 UTC")
	expectedTime2, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 10:55:00 +0000 UTC")
	expected := []time.Time{expectedTime1, expectedTime2}

	// Assert
	if reflect.DeepEqual(actual, expected) {
		t.Error(actual)
	}
}

func Test_memory_graph_data_should_not_be_nil(t *testing.T) {
	//Arrange
	test := ByMemory{
		DocCountErrorUpperBound: 0,
		SumOtherDocCount:        0,
		Buckets:                 []MemoryBuckets{{
			Key:                "",
			DocCount:           0,
			GraphDataForMemory: GraphDataForMemory{Buckets: []ForMemoryBuckets{{
				KeyAsString:   time.Time{},
				Key:           0,
				DocCount:      0,
				MemoryAverage: MemoryAverage{65},
			},{
				KeyAsString:   time.Time{},
				Key:           0,
				DocCount:      0,
				MemoryAverage: MemoryAverage{30},
			}},
				},
		}},
	}
	expected := []float64{65, 30}
	//Act
	actual := test.GetMemoryGraphData(100)

	// Assert
	fmt.Println(actual)
	fmt.Println(expected)
	if reflect.DeepEqual(actual, expected)  == false {
		t.Error(actual)
	}
}

func Test_disk_graph_data_should_not_be_nil(t *testing.T) {
	//Arrange
	test := ByDisk{
		DocCountErrorUpperBound: 0,
		SumOtherDocCount:        0,
		Buckets:                 []DiskBuckets{{
			Key:                   "",
			DocCount:              0,
			GraphDataForDiskSpace: GraphDataForDiskSpace{Buckets:[]DiskSpaceBuckets{{
				KeyAsString: time.Time{},
				Key:         0,
				DocCount:    0,
				DiskAverage: DiskAverage{65},
			},{
				KeyAsString: time.Time{},
				Key:         0,
				DocCount:    0,
				DiskAverage: DiskAverage{25},
			},
			}},
		}},
	}
	expected := []float64{65, 25}
	//Act
	actual := test.GetDiskGraphData(100)

	// Assert
	if reflect.DeepEqual(actual, expected)  == false {
		t.Error(actual)
	}
}
