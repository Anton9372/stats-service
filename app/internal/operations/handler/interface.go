package handler

import (
	"context"
	"stats-service/internal/operations/entity"
	"stats-service/pkg/api/sort"
)

type Service interface {
	GetAll(ctx context.Context, sortOptions sort.Options) ([]entity.Operation, error)
}
