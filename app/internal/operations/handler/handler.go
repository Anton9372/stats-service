package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"stats-service/internal/apperror"
	h "stats-service/internal/handler"
	"stats-service/pkg/api/sort"
	"stats-service/pkg/logging"
	"stats-service/pkg/utils"
)

const (
	operationsURL = "/api/operations"
	//userByIdURL = "/api/users/one/:uuid"
	//allUsersURL = "/api/users/all"
)

type handler struct {
	service Service
	logger  *logging.Logger
}

func NewHandler(service Service, logger *logging.Logger) h.Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, operationsURL, sort.Middleware(apperror.Middleware(h.GetOperations), "date_time", sort.ASC))
}

func (h *handler) GetOperations(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("get ops")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	var sortOptions sort.Options
	if options, ok := r.Context().Value(sort.OptionsContextKey).(sort.Options); ok {
		sortOptions = options
	}

	operations, err := h.service.GetAll(r.Context(), sortOptions)
	if err != nil {
		return err
	}

	dataBytes, err := json.Marshal(operations)
	if err != nil {
		return fmt.Errorf("failed to marshal operations: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(dataBytes)
	h.logger.Info("success")
	return nil
}
