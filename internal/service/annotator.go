package service

import (
	"context"
	"time"

	"github.com/blackmou5e/grafana-annotator/internal/grafana"
	"github.com/blackmou5e/grafana-annotator/internal/validation"

	"github.com/sirupsen/logrus"
)

type AnnotatorService struct {
	client    grafana.GrafanaClient
	logger    *logrus.Logger
	validator *validation.Validator
}

func NewAnnotatorService(client grafana.GrafanaClient, logger *logrus.Logger, validator *validation.Validator) *AnnotatorService {
	return &AnnotatorService{
		client:    client,
		logger:    logger,
		validator: validator,
	}
}

func (s *AnnotatorService) CreateAnnotations(ctx context.Context, tags []string, message string) error {
	if err := s.validator.ValidateAnnotationInput(tags, message); err != nil {
		return err
	}

	dashboards, err := s.client.FetchDashboards(ctx)
	if err != nil {
		return err
	}

	s.logger.WithField("dashboard_count", len(dashboards)).Debug("Fetched dashboards")

	for _, dashboard := range dashboards {
		annotation := grafana.Annotation{
			DashboardUID: dashboard.UID,
			TimeStart:    time.Now().UTC().UnixMilli(),
			Tags:         tags,
			Text:         message,
		}

		if err := s.client.CreateAnnotation(ctx, annotation); err != nil {
			s.logger.WithError(err).WithField("dashboard", dashboard.Title).
				Error("Failed to create annotation")
			continue
		}

		s.logger.WithField("dashboard", dashboard.Title).Info("Created annotation")
	}

	return nil
}
