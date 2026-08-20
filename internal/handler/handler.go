package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	models "github.com/cyntraten/metrics-ed/internal/model"
	"github.com/cyntraten/metrics-ed/internal/repository"
)

type Handler struct {
	storage repository.Storage
}

func getQueryParts(urlPath string) (parts []string, error error, errCode int) {
	parts = strings.Split(urlPath, "/")

	// 404 должен возвращаться только если нет имени метрики, остальное 400
	switch len(parts) {
	case 4:
		metricName := parts[3]
		if len(metricName) == 0 {
			return []string{}, fmt.Errorf("Empty metric name"), http.StatusNotFound
		} else {
			return []string{}, fmt.Errorf("Empty metric value"), http.StatusBadRequest
		}
	case 5:
		return parts, nil, 0
	default:
		return []string{}, fmt.Errorf("Data is not valid"), http.StatusBadRequest
	}
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

func (h *Handler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	// Вынести в отдельную функцию getQueryParts
	// Вынес!
	parts, err, errCode := getQueryParts(r.URL.Path)
	if err != nil {
		switch errCode {
		case 404:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case 400:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
	case string(models.Gauge):
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Metric value is not valid", http.StatusBadRequest)
			return
		}

		var GaugeMetric models.Metrics
		GaugeMetric.ID = metricName
		GaugeMetric.MType = models.Gauge
		GaugeMetric.Value = &value

		err = h.storage.UpdateMetric(metricName, GaugeMetric)

		if err != nil {
			http.Error(w, "Save metric error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		return
	case string(models.Counter):
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Metric value is not valid", http.StatusBadRequest)
			return
		}

		var CounterMetric models.Metrics
		CounterMetric.ID = metricName
		CounterMetric.MType = models.Counter
		CounterMetric.Delta = &value

		err = h.storage.UpdateMetric(metricName, CounterMetric)

		if err != nil {
			http.Error(w, "Save metric error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		return
	default:
		http.Error(w, "Metric type is not valid", http.StatusBadRequest)
		return

	}

}

func NewHandler(storage repository.Storage) *Handler {
	h := Handler{storage: storage}
	return &h
}
