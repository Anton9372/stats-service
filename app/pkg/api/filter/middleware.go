package filter

import (
	"context"
	"net/http"
	"strconv"
)

const (
	OptionsContextKey = "filter_options"
)

func Middleware(h http.HandlerFunc, defaultLimit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitFromQuery := r.URL.Query().Get("limit")

		limit := defaultLimit
		var limitParseErr error
		if limitFromQuery != "" {
			if limit, limitParseErr = strconv.Atoi(limitFromQuery); limitParseErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("invalid limit"))
				return
			}
		}

		optionsWithLimit := NewOptions(limit)
		ctx := context.WithValue(r.Context(), OptionsContextKey, optionsWithLimit)
		r = r.WithContext(ctx)

		h(w, r)
	}
}
