package storage

import (
	"math/rand"
	"reflect"
	"runtime"

	"github.com/anatolev-max/metrics-tpl/internal/enum"
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
	UpdateMetricValue(name string, value any)
}

type MemStorage struct {
	Counter map[string]int64
	Gauge   map[string]float64
}

func NewMemStorage() MemStorage {
	return MemStorage{
		Counter: map[string]int64{
			enum.PollCounter: 0,
		},
		Gauge: map[string]float64{
			enum.RandomValue: 0,
		},
	}
}

func (s *MemStorage) UpdateAgentData() {
	rtm := runtime.MemStats{}
	runtime.ReadMemStats(&rtm)

	gauges := Metrics{}.Gauges
	gaugesValue := reflect.ValueOf(gauges)

	for fieldIndex := 0; fieldIndex < gaugesValue.NumField(); fieldIndex++ {
		name := gaugesValue.Type().Field(fieldIndex).Name
		value := reflect.ValueOf(rtm).FieldByName(name)

		if value.IsValid() {
			switch vType := value.Interface().(type) {
			case uint64:
				s.Gauge[name] = float64(vType)
			case uint32:
				s.Gauge[name] = float64(vType)
			case float64:
				s.Gauge[name] = vType
			}

			s.Counter[enum.PollCounter]++
		}
	}

	s.Gauge[enum.RandomValue] = rand.Float64()
}

func (s *MemStorage) UpdateMetricValue(name string, value any) {
	switch vType := value.(type) {
	case int64:
		if _, exist := s.Counter[name]; !exist || name == enum.PollCounter {
			s.Counter[name] = vType
		} else {
			s.Counter[name] += vType
		}
	case float64:
		s.Gauge[name] = vType
	}
}
