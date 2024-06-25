package service

import (
	"context"
	"fmt"
	"stats-service/internal/operations/entity"
	"stats-service/internal/operations/handler"
	"stats-service/internal/operations/storage"
	"stats-service/pkg/api/sort"
	"stats-service/pkg/logging"
)

type service struct {
	repository storage.Repository
	logger     *logging.Logger
}

func NewService(repository storage.Repository, logger *logging.Logger) handler.Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) GetAll(ctx context.Context, sortOptions sort.Options) ([]entity.Operation, error) {
	options := storage.NewSortOptions(sortOptions.Field, sortOptions.Order)

	operations, err := s.repository.FindAll(ctx, options)
	if err != nil {
		return operations, fmt.Errorf("failed to get operations: %w", err)
	}
	return operations, nil
}
