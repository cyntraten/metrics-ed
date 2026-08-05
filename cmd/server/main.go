package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	Type  MetricType
	Value float64
}

type MemStorage struct {
	Metrics map[string]Metric
}

type Storage interface {
	UpdateMetric(name string, metric Metric)
	GetMetric(name string) (Metric, bool)
}

func getQueryParts(urlPath string) (parts []string, error error) {
	parts = strings.Split(urlPath, "/")
	if len(parts) != 5 {
		return []string{}, fmt.Errorf("Invalid query params")
	}

	return parts, nil
}

func getMetricsDataFromParts(parts []string) (metricType string, metricName string, metricValue string) {
	metricType = strings.ToLower(parts[2])
	metricName = parts[3]
	metricValue = parts[4]

	return metricType, metricName, metricValue
}

func validateMetricName(metricName string) (err error) {
	if utf8.RuneCountInString(metricName) <= 0 {
		return fmt.Errorf("Metric name is empty")
	}

	return nil
}

func updateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	// Вынести в отдельную функцию getQueryParts
	// Вынес!
	parts, err := getQueryParts(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вынес функцию для получения данных метрик из частей запроса
	metricType, metricName, metricValue := getMetricsDataFromParts(parts)

	// Проверяем типы данных в зависимости
	// от типа метрики в идеале тоже вынести
	// в другую функцию (вынес)

	//validateMetricName

	err = validateMetricName(metricName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Поменял if на switch
	switch metricType {
	case string(Gauge):
		_, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Metric value is not valid", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case string(Counter):
		_, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Metric value is not valid", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	default:
		http.Error(w, "Metric type is not valid", http.StatusBadRequest)
		return

	}

}

func main() {

	//savedMetrics := MemStorage{}

	mux := http.NewServeMux()

	mux.HandleFunc(`/update/`, updateMetrics)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
