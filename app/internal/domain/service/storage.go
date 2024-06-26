package service

import (
	"context"
	"stats-service/internal/domain/entity"
	"stats-service/internal/storage/sorting"
	"stats-service/pkg/api/filter"
)

type Repository interface {
	FindAll(ctx context.Context, sortOptions sorting.SortOptions, filterOptions filter.Options) ([]entity.Operation, error)
}
