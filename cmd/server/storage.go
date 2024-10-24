package main

type MemStorageInterface interface {
	Update(name string, value interface{})
}

type MemStorage struct {
	Counter map[string]int64
	Gauge   map[string]float64
}

func NewMemStorage() MemStorage {
	return MemStorage{
		Counter: map[string]int64{},
		Gauge:   map[string]float64{},
	}
}

func (ms *MemStorage) Update(name string, value interface{}) {
	switch vv := value.(type) {
	case int64:
		if _, exist := ms.Counter[name]; !exist {
			ms.Counter[name] = vv
		} else {
			ms.Counter[name] += vv
		}
	case float64:
		ms.Gauge[name] = vv
	}
}
