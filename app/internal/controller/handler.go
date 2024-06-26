package controller

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"stats-service/internal/apperror"
	"stats-service/internal/domain/entity"
	"stats-service/pkg/api/filter"
	"stats-service/pkg/api/sort"
	"stats-service/pkg/logging"
	"stats-service/pkg/utils"
	"strconv"
	"strings"
	"time"
)

const (
	operationsURL = "/api/operations"
)

type handler struct {
	service Service
	logger  *logging.Logger
}

func NewHandler(service Service, logger *logging.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, operationsURL,
		filter.Middleware(sort.Middleware(apperror.Middleware(h.GetOperations), entity.DateTime, sort.ASC), 20))
}

func (h *handler) GetOperations(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Get operations")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	var sortOptions sort.Options
	if options, ok := r.Context().Value(sort.OptionsContextKey).(sort.Options); ok {
		sortOptions = options
	}

	filterOptions := r.Context().Value(filter.OptionsContextKey).(filter.Options)

	categoryUUID := r.URL.Query().Get(entity.CategoryUUID)
	if categoryUUID != "" {
		err := filterOptions.AddField(entity.CategoryUUID, filter.OperatorEqual, categoryUUID, filter.DataTypeFloat)
		if err != nil {
			return apperror.BadRequestError(err.Error())
		}
	}

	moneySum := r.URL.Query().Get(entity.MoneySum)
	if moneySum != "" {
		operator := filter.OperatorEqual
		value := moneySum

		if strings.Index(moneySum, ":") != -1 {
			split := strings.Split(moneySum, ":")
			operator = split[0]
			value = split[1]
		}

		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			validationErr := apperror.BadRequestError("filter params validation failed")
			validationErr.WithParams(map[string]string{
				"money_sum": "this field should be a number",
			})
			return validationErr
		}

		err = filterOptions.AddField(entity.MoneySum, operator, value, filter.DataTypeFloat)
		if err != nil {
			return apperror.BadRequestError(err.Error())
		}
	}

	dateTime := r.URL.Query().Get(entity.DateTime)
	if dateTime != "" {
		operator := filter.OperatorBetween
		var err error

		if strings.Index(dateTime, ":") != -1 {
			split := strings.Split(dateTime, ":")
			dateBegin := split[0]
			dateEnd := split[1]

			_, err = time.Parse(time.DateOnly, dateBegin)
			_, err = time.Parse(time.DateOnly, dateEnd)
		} else {
			_, err = time.Parse(time.DateOnly, dateTime)
			dateTime = fmt.Sprintf("%s:%s", dateTime, dateTime)
		}

		if err != nil {
			validationErr := apperror.BadRequestError("filter params validation failed")
			validationErr.WithParams(map[string]string{
				"date_time": fmt.Sprintf("date should be in format: %s", time.DateOnly),
			})
			return validationErr
		}

		err = filterOptions.AddField(entity.DateTime, operator, dateTime, filter.DataTypeDate)
		if err != nil {
			return apperror.BadRequestError(err.Error())
		}
	}

	report, err := h.service.GetAll(r.Context(), sortOptions, filterOptions)
	if err != nil {
		return err
	}

	dataBytes, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal operations: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(dataBytes)
	h.logger.Info("Get operations successfully")
	return nil
}
