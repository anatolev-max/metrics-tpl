package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/anatolev-max/metrics-tpl/internal/config"
	"github.com/anatolev-max/metrics-tpl/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetWebHook(t *testing.T) {
	testCases := []struct {
		name         string
		method       string
		expectedCode int
		contentType  string
		metricType   string
		metricName   string
		metricValue  interface{}
	}{
		{
			name:         "test #1 - ok Counter",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			contentType:  config.TextPlain,
			metricType:   config.Counter,
			metricName:   "Go",
			metricValue:  123,
		},
		{
			name:         "test #2 - ok Gauge",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			contentType:  config.TextPlain,
			metricType:   config.Gauge,
			metricName:   "Go",
			metricValue:  123.1,
		},
		{
			name:         "test #3 - unsupported method",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			contentType:  config.TextPlain,
			metricType:   config.Counter,
			metricName:   "Go",
			metricValue:  123,
		},
		{
			name:         "test #4 - unsupported Content-Type",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  config.ApplicationJson,
			metricType:   config.Counter,
			metricName:   "Go",
			metricValue:  123,
		},
		{
			name:         "test #5 - without metricName",
			method:       http.MethodPost,
			expectedCode: http.StatusNotFound,
			contentType:  config.TextPlain,
			metricType:   config.Counter,
			metricName:   "",
			metricValue:  123,
		},
		{
			name:         "test #6 - unsupported metricType",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  config.TextPlain,
			metricType:   "Golang",
			metricName:   "Go",
			metricValue:  123,
		},
		{
			name:         "test #6 - unsupported metricValue fot Counter",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  config.TextPlain,
			metricType:   config.Counter,
			metricName:   "Go",
			metricValue:  123.1,
		},
		{
			name:         "test #7 - unsupported metricValue fot Gauge",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			contentType:  config.TextPlain,
			metricType:   config.Gauge,
			metricName:   "Go",
			metricValue:  "Golang",
		},
	}

	for _, tc := range testCases {
		s := storage.NewMemStorage()

		t.Run(tc.name, func(t *testing.T) {
			var url string
			urlPattern := config.UpdateFullEndpoint + "%v/%v/%v"

			if reflect.ValueOf(tc.metricValue).Kind() == reflect.Int {
				url = fmt.Sprintf(urlPattern, tc.metricType, tc.metricName, int64(tc.metricValue.(int)))
			} else {
				url = fmt.Sprintf(urlPattern, tc.metricType, tc.metricName, tc.metricValue)
			}

			request := httptest.NewRequest(tc.method, url, nil)
			request.Header.Add("Content-Type", tc.contentType)
			writer := httptest.NewRecorder()

			handler := GetUpdateWebhook(s)
			handler(writer, request)

			assert.Equal(t, tc.expectedCode, writer.Code, "The response code does not match what is expected")
		})
	}
}