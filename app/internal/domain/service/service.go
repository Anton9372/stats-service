package service

import (
	"context"
	"fmt"
	"stats-service/internal/controller"
	"stats-service/internal/domain/entity"
	"stats-service/internal/storage/sorting"
	"stats-service/pkg/api/filter"
	"stats-service/pkg/api/sort"
	"stats-service/pkg/logging"
)

type service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) controller.Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) GetAll(ctx context.Context, sortOptions sort.Options, filterOptions filter.Options) (entity.Report, error) {
	var report entity.Report
	sortOpt := sorting.NewSortOptions(sortOptions.Field, sortOptions.Order)

	operations, err := s.repository.FindAll(ctx, sortOpt, filterOptions)
	if err != nil {
		return report, fmt.Errorf("failed to get operations: %w", err)
	}

	report = entity.NewReport(operations)
	return report, nil
}
