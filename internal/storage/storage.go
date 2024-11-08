package storage

import (
	"math/rand"
	"reflect"
	"runtime"

	"github.com/anatolev-max/metrics-tpl/internal/config"
)

type Metrics struct {
	Gauges struct {
		Alloc,
		BuckHashSys,
		Frees,
		GCSys,
		HeapAlloc,
		HeapIdle,
		HeapInuse,
		HeapObjects,
		HeapReleased,
		HeapSys,
		LastGC,
		Lookups,
		MCacheInuse,
		MCacheSys,
		MSpanInuse,
		MSpanSys,
		Mallocs,
		NextGC,
		OtherSys,
		PauseTotalNs,
		StackInuse,
		StackSys,
		Sys,
		TotalAlloc uint64

		GCCPUFraction float64

		NumForcedGC uint32
		NumGC       uint32
	}
}

type Storage interface {
	UpdateAgentData()
	UpdateServerData(name string, value interface{})
}

type MemStorage struct {
	Counter map[string]int64
	Gauge   map[string]float64
}

func NewMemStorage() MemStorage {
	return MemStorage{
		Counter: map[string]int64{
			config.PollCounter: 0,
		},
		Gauge: map[string]float64{
			config.RandomValue: 0,
		},
	}
}

func (ms *MemStorage) UpdateAgentData() {
	rtm := runtime.MemStats{}
	runtime.ReadMemStats(&rtm)

	gauges := Metrics{}.Gauges
	val := reflect.ValueOf(gauges)

	if val.Kind() == reflect.Struct {
		for fieldIndex := 0; fieldIndex < val.NumField(); fieldIndex++ {
			name := val.Type().Field(fieldIndex).Name
			value := reflect.ValueOf(rtm).FieldByName(name)

			if value.IsValid() {
				switch vv := value.Interface().(type) {
				case uint64:
					ms.Gauge[name] = float64(vv)
				case uint32:
					ms.Gauge[name] = float64(vv)
				case float64:
					ms.Gauge[name] = vv
				}

				ms.Counter[config.PollCounter]++
			}
		}

		ms.Gauge[config.RandomValue] = rand.Float64()
	}
}

func (ms *MemStorage) UpdateServerData(name string, value interface{}) {
	switch vType := value.(type) {
	case int64:
		if _, exist := ms.Counter[name]; !exist || name == config.PollCounter {
			ms.Counter[name] = vType
		} else {
			ms.Counter[name] += vType
		}
	case float64:
		ms.Gauge[name] = vType
	}
}
