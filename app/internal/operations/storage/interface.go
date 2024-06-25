package storage

import (
	"context"
	"stats-service/internal/operations/entity"
)

type Repository interface {
	FindAll(ctx context.Context, sortOptions SortOptions) ([]entity.Operation, error)
}

type SortOptions interface {
	GetOrderBy() string
}
