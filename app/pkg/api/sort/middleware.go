package sort

import (
	"context"
	"net/http"
)

const (
	ASC               = "ASC"
	DESC              = "DESC"
	OptionsContextKey = "sort_options"
)

func Middleware(h http.HandlerFunc, defaultSortField, defaultSortOrder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := r.URL.Query().Get("sort_order")

		if sortBy == "" {
			sortBy = defaultSortField
		}

		if sortOrder == "" {
			sortOrder = defaultSortOrder
		}

		options := Options{
			Field: sortBy,
			Order: sortOrder,
		}

		ctx := context.WithValue(r.Context(), OptionsContextKey, options)
		r = r.WithContext(ctx)

		h(w, r)
	}
}
