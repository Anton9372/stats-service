package sorting

import (
	"fmt"
	"stats-service/internal/apperror"
	"stats-service/internal/domain/entity"
	"stats-service/pkg/api/sort"
	"strings"
)

type sortOptions struct {
	Field, Order string
}

func NewSortOptions(field, order string) (SortOptions, error) {
	if err := validateField(field); err != nil {
		return nil, err
	}
	if err := validateOrder(order); err != nil {
		return nil, err
	}
	return &sortOptions{
		Field: field,
		Order: order,
	}, nil
}

func (so *sortOptions) GetOrderBy() string {
	return fmt.Sprintf("%s %s", so.Field, so.Order)
}

func validateField(field string) error {
	switch field {
	case entity.MoneySum:
	case entity.Description:
	case entity.DateTime:
	default:
		err := apperror.BadRequestError("sort field validation failed")
		err.WithFields(map[string]string{
			"sort_by": fmt.Sprintf("possible fields: %s, %s, %s",
				entity.MoneySum, entity.Description, entity.DateTime),
		})
		return err
	}
	return nil
}

func validateOrder(order string) error {
	upperSortOrder := strings.ToUpper(order)
	if upperSortOrder != sort.ASC && upperSortOrder != sort.DESC {
		err := apperror.BadRequestError("sort order validation failed")
		err.WithFields(map[string]string{
			"sort_order": fmt.Sprintf("possible orders: %s, %s",
				sort.ASC, sort.DESC),
		})
		return err
	}
	return nil
}
