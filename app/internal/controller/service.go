package controller

import (
	"context"
	"stats-service/internal/domain/entity"
	"stats-service/pkg/api/filter"
	"stats-service/pkg/api/sort"
)

type Service interface {
	GetAll(ctx context.Context, sortOptions sort.Options, filterOptions filter.Options) (entity.Report, error)
}
