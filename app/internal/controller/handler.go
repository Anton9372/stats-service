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
	"strings"
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

	validationErr := apperror.BadRequestError("filter params validation failed")

	userUUID := r.URL.Query().Get(entity.UserUUID)
	if userUUID != "" {
		err := filterOptions.AddField(entity.UserUUID, filter.OperatorEqual, userUUID, filter.DataTypeString)
		if err != nil {
			validationErr.WithParams(map[string]string{
				entity.UserUUID: err.Error(),
			})
			return validationErr
		}
	}

	categoryName := r.URL.Query().Get(entity.CategoryName)
	if categoryName != "" {
		err := filterOptions.AddField(entity.CategoryName, filter.OperatorSubString, categoryName, filter.DataTypeString)
		if err != nil {
			validationErr.WithParams(map[string]string{
				entity.CategoryName: err.Error(),
			})
			return validationErr
		}
	}

	categoryType := r.URL.Query().Get(entity.TypeOfCategory)
	if categoryType != "" {
		err := filterOptions.AddField(entity.TypeOfCategory, filter.OperatorEqual, categoryType, filter.DataTypeString)
		if err != nil {
			validationErr.WithParams(map[string]string{
				entity.TypeOfCategory: err.Error(),
			})
			return validationErr
		}
	}

	categoryUUID := r.URL.Query().Get(entity.CategoryUUID)
	if categoryUUID != "" {
		err := filterOptions.AddField(entity.CategoryUUID, filter.OperatorEqual, categoryUUID, filter.DataTypeString)
		if err != nil {
			validationErr.WithParams(map[string]string{
				entity.CategoryUUID: err.Error(),
			})
			return validationErr
		}
	}

	description := r.URL.Query().Get(entity.Description)
	if description != "" {
		err := filterOptions.AddField(entity.Description, filter.OperatorSubString, description, filter.DataTypeString)
		if err != nil {
			validationErr.WithParams(map[string]string{
				entity.Description: err.Error(),
			})
			return validationErr
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

		err := filterOptions.AddField(entity.MoneySum, operator, value, filter.DataTypeFloat)
		if err != nil {
			validationErr.WithParams(map[string]string{
				entity.MoneySum: err.Error(),
			})
			return validationErr
		}
	}

	dateTime := r.URL.Query().Get(entity.DateTime)
	if dateTime != "" {
		operator := filter.OperatorBetween

		if strings.Index(dateTime, ":") == -1 {
			dateTime = fmt.Sprintf("%s:%s", dateTime, dateTime)
		}

		err := filterOptions.AddField(entity.DateTime, operator, dateTime, filter.DataTypeDate)
		if err != nil {
			validationErr.WithParams(map[string]string{
				entity.DateTime: err.Error(),
			})
			return validationErr
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
