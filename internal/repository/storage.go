package repository

import (
	"fmt"

	models "github.com/cyntraten/metrics-ed/internal/model"
)

type MemStorage struct {
	Metrics map[string]models.Metrics
}

type Storage interface {
	UpdateMetric(name string, metrics models.Metrics) error
	GetMetric(name string) (models.Metrics, bool)
}

func getMType(metricType string) (string, error) {
	if metricType == models.Gauge {
		return models.Gauge, nil
	}
	if metricType == models.Counter {
		return models.Counter, nil
	}

	return "", fmt.Errorf("Unexpected MType")
}

func (s *MemStorage) UpdateMetric(name string, metrics models.Metrics) error {
	value, ok := s.Metrics[name]

	metricType, err := getMType(metrics.MType)
	if err != nil {
		return err
	}

	if metricType == models.Gauge {
		if metrics.Value != nil {
			if !ok {
				s.Metrics[name] = metrics
				return nil
			} else {
				*s.Metrics[name].Value = *metrics.Value
				return nil
			}
		} else {
			return fmt.Errorf("Metrics: %s has not value", metricType)
		}
	}

	if metricType == models.Counter {
		if metrics.Delta != nil {
			if !ok {
				s.Metrics[name] = metrics
				return nil
			} else {
				*s.Metrics[name].Delta = *metrics.Delta + *value.Delta
				return nil
			}
		} else {
			return fmt.Errorf("Metrics: %s has not value", metricType)
		}
	}

	// Пока не могу понять что недопроверил и почему выдается ошибка компилятором.
	return fmt.Errorf("Unexpected error in UpdateMetric")

}

func (s *MemStorage) GetMetric(name string) (models.Metrics, bool) {
	value, ok := s.Metrics[name]
	if !ok {
		return value, false
	}
	return value, true
}

func NewMemStorage() *MemStorage {
	s := MemStorage{
		Metrics: make(map[string]models.Metrics),
	}
	return &s
}
