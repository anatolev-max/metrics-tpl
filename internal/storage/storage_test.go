package storage

import (
	"reflect"
	"testing"

	"github.com/anatolev-max/metrics-tpl/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateAgentData(t *testing.T) {
	const myGaugeCount int = 1

	s := NewMemStorage()
	assert.Equal(t, len(s.Gauge), myGaugeCount)
	gaugeCount := reflect.ValueOf(Metrics{}.Gauges).NumField()

	s.UpdateAgentData()
	storageGaugeCount := len(s.Gauge)

	assert.Equal(t, storageGaugeCount, gaugeCount+myGaugeCount)
	assert.Equal(t, int(s.Counter[config.PollCounter]), gaugeCount)
}

func TestMemStorage_UpdateMetricValue(t *testing.T) {
	testCases := []struct {
		name        string
		metricName  string
		metricValue interface{}
	}{
		{
			name:        "test #1 - ok Counter",
			metricName:  "Golang",
			metricValue: 123,
		},
		{
			name:        "test #2 - ok Gauge",
			metricName:  "Golang",
			metricValue: 123.1,
		},
	}

	for _, tc := range testCases {
		const updateCount int = 3
		s := NewMemStorage()

		t.Run(tc.name, func(t *testing.T) {
			switch reflect.ValueOf(tc.metricValue).Kind() {
			case reflect.Int:
				convertedValue := int64(tc.metricValue.(int))
				for i := 1; i <= updateCount; i++ {
					s.UpdateMetricValue(tc.metricName, convertedValue)
					assert.Equal(t, convertedValue*int64(i), s.Counter[tc.metricName])
				}
			case reflect.Float64:
				for i := 0; i < updateCount; i++ {
					convertedValue := tc.metricValue.(float64) + float64(i)
					s.UpdateMetricValue(tc.metricName, convertedValue)
					assert.Equal(t, convertedValue, s.Gauge[tc.metricName])
				}
			default:
			}
		})
	}
}
