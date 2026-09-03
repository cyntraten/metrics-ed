package repository

import (
	"testing"

	models "github.com/cyntraten/metrics-ed/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricGauge(t *testing.T) {

	var valueEmptyStorage float64 = 21.5
	var valueForNonEmptyStorage float64 = 20.1

	tests := []struct {
		name           string
		initialStorage *models.Metrics
		metricToUpdate models.Metrics
		expectedMetric models.Metrics
		expectedValue  float64
		wantErr        bool
	}{
		{
			name:           "Add gauge to empty storage",
			initialStorage: nil,
			metricToUpdate: models.Metrics{
				ID:    "TestMetric1",
				MType: models.Gauge,
				Value: &valueEmptyStorage,
			},
			expectedMetric: models.Metrics{
				ID:    "TestMetric1",
				MType: models.Gauge,
				Value: &valueEmptyStorage,
			},
		},
		{
			name: "Add gauge to non empty storage",
			initialStorage: &models.Metrics{
				ID:    "TestMetric2",
				MType: models.Gauge,
				Value: &valueEmptyStorage,
			},
			metricToUpdate: models.Metrics{
				ID:    "TestMetric2",
				MType: models.Gauge,
				Value: &valueForNonEmptyStorage,
			},
			expectedMetric: models.Metrics{
				ID:    "TestMetric2",
				MType: models.Gauge,
				Value: &valueForNonEmptyStorage,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			storage := NewMemStorage()

			if test.initialStorage != nil {
				storage.UpdateMetric(test.initialStorage.ID, *test.initialStorage)
			}

			err := storage.UpdateMetric(test.metricToUpdate.ID, test.metricToUpdate)
			assert.NoError(t, err)
			storageAnswer, ok := storage.GetMetric(test.metricToUpdate.ID)

			if !ok {
				t.Errorf("Metric not added to empty storage")
			}

			assert.Equal(t, test.expectedMetric, storageAnswer)
		})

	}

}

func TestUpdateMetricCounter(t *testing.T) {

	var valueEmptyStorage int64 = 21
	var valueForNonEmptyStorage int64 = 20
	var valueAfterUpdateStorage int64 = valueEmptyStorage + valueForNonEmptyStorage

	tests := []struct {
		name           string
		initialStorage *models.Metrics
		metricToUpdate models.Metrics
		expectedMetric models.Metrics
		expectedValue  int
		wantErr        bool
	}{
		{
			name:           "Add counter to empty storage",
			initialStorage: nil,
			metricToUpdate: models.Metrics{
				ID:    "TestMetric1",
				MType: models.Counter,
				Delta: &valueEmptyStorage,
			},
			expectedMetric: models.Metrics{
				ID:    "TestMetric1",
				MType: models.Counter,
				Delta: &valueEmptyStorage,
			},
		},
		{
			name: "Add counter to non empty storage",
			initialStorage: &models.Metrics{
				ID:    "TestMetric2",
				MType: models.Counter,
				Delta: &valueEmptyStorage,
			},
			metricToUpdate: models.Metrics{
				ID:    "TestMetric2",
				MType: models.Counter,
				Delta: &valueForNonEmptyStorage,
			},
			expectedMetric: models.Metrics{
				ID:    "TestMetric2",
				MType: models.Counter,
				Delta: &valueAfterUpdateStorage,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			storage := NewMemStorage()

			if test.initialStorage != nil {
				storage.UpdateMetric(test.initialStorage.ID, *test.initialStorage)
			}

			err := storage.UpdateMetric(test.metricToUpdate.ID, test.metricToUpdate)
			assert.NoError(t, err)
			storageAnswer, ok := storage.GetMetric(test.metricToUpdate.ID)

			if !ok {
				t.Errorf("Metric not added to empty storage")
			}

			assert.Equal(t, test.expectedMetric, storageAnswer)
		})

	}

}
