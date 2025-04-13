package tests

import (
	"context"
	"testing"

	"github.com/blackmou5e/grafana-annotator/internal/grafana"
	"github.com/blackmou5e/grafana-annotator/internal/service"
	"github.com/blackmou5e/grafana-annotator/internal/validation"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGrafanaClient struct {
	mock.Mock
}

var _ grafana.GrafanaClient = (*MockGrafanaClient)(nil)

func (m *MockGrafanaClient) FetchDashboards(ctx context.Context) ([]grafana.Dashboard, error) {
	args := m.Called(ctx)
	return args.Get(0).([]grafana.Dashboard), args.Error(1)
}

func (m *MockGrafanaClient) CreateAnnotation(ctx context.Context, annotation grafana.Annotation) error {
	args := m.Called(ctx, annotation)
	return args.Error(0)
}

func TestAnnotatorService(t *testing.T) {
	mockClient := new(MockGrafanaClient)
	logger := logrus.New()
	validator := validation.NewValidator()

	service := service.NewAnnotatorService(mockClient, logger, validator)

	dashboards := []grafana.Dashboard{
		{
			UID:   "test-dashboard",
			Title: "Test Dashboard",
		},
	}

	ctx := context.Background()
	tags := []string{"test"}
	message := "Test annotation"

	mockClient.On("FetchDashboards", ctx).Return(dashboards, nil)
	mockClient.On("CreateAnnotation", ctx, mock.Anything).Return(nil)

	err := service.CreateAnnotations(ctx, tags, message)
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}
