package storage

import (
	"reflect"
	"testing"

	"github.com/anatolev-max/metrics-tpl/cmd/common"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateAgentData(t *testing.T) {
	const myGaugeCount int = 1

	ms := NewMemStorage()
	assert.Equal(t, len(ms.Gauge), myGaugeCount)
	gaugeCount := reflect.ValueOf(Metrics{}.Gauges).NumField()

	ms.UpdateAgentData()
	storageGaugeCount := len(ms.Gauge)

	assert.Equal(t, storageGaugeCount, gaugeCount+myGaugeCount)
	assert.Equal(t, int(ms.Counter[common.PollCounter]), gaugeCount)
}

func TestMemStorage_UpdateServerData(t *testing.T) {
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
		updateCount := 3
		ms := NewMemStorage()

		t.Run(tc.name, func(t *testing.T) {
			switch reflect.ValueOf(tc.metricValue).Kind() {
			case reflect.Int:
				convertedValue := int64(tc.metricValue.(int))
				for i := 1; i <= updateCount; i++ {
					ms.UpdateServerData(tc.metricName, convertedValue)
					assert.Equal(t, convertedValue*int64(i), ms.Counter[tc.metricName])
				}
			case reflect.Float64:
				for i := 0; i < updateCount; i++ {
					convertedValue := tc.metricValue.(float64) + float64(i)
					ms.UpdateServerData(tc.metricName, convertedValue)
					assert.Equal(t, convertedValue, ms.Gauge[tc.metricName])
				}
			default:
			}
		})
	}
}
