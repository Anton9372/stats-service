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
	operationsURL = "/api/stats"
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

// GetOperations
// @Summary 	Get operations
// @Description Retrieves a list of operations with support for filtering and sorting.
// @Tags 		Operations
// @Id 			get-operations
// @Produce 	json
// @Param 		user_uuid 	  path 	   string false  "User UUID"
// @Param 		category_name path 	   string false  "Category name"
// @Param 		type query 	  path 	   string false  "Category type"
// @Param 		category_id   path 	   string false  "Category ID"
// @Param 		description   path 	   string false  "Description"
// @Param 		money_sum 	  path 	   string false  "Money sum (supports operators for numbers: eq, neq, lt, lte, gt, gte)"
// @Param 		date_time     path 	   string false  "Date and time of operation (supports formats: yyyy-mm-dd, yyyy-mm-dd:yyyy-mm-dd)"
// @Param 		sort_by 	  path 	   string false  "Field to sort by (money_sum, date_time, description)"
// @Param 		sort_order 	  path 	   string false  "Sort order (asc, desc)"
// @Success 	200 		  {object} string 		 "List of operations"
// @Router /operations [get]
func (h *handler) GetOperations(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Get operations")
	defer utils.CloseBody(h.logger, r.Body)
	w.Header().Set("Content-Type", "application/json")

	var sortOptions sort.Options
	if options, ok := r.Context().Value(sort.OptionsContextKey).(sort.Options); ok {
		sortOptions = options
	}
	// @Failure 	400 		  {object} apperror "Validation error in filter or sort parameters"
	// @Failure 	418 		  {object} apperror "Something wrong with application logic"
	// @Failure 	500 		  {object} apperror "Internal server error"
	filterOptions := r.Context().Value(filter.OptionsContextKey).(filter.Options)

	var err error
	userUUID := r.URL.Query().Get(entity.UserUUID)
	filterOptions, err = processParam(userUUID, filter.DataTypeString, entity.UserUUID, filterOptions)
	if err != nil {
		return err
	}

	categoryName := r.URL.Query().Get(entity.CategoryName)
	filterOptions, err = processParam(categoryName, filter.DataTypeString, entity.CategoryName, filterOptions)
	if err != nil {
		return err
	}

	categoryType := r.URL.Query().Get(entity.TypeOfCategory)
	filterOptions, err = processParam(categoryType, filter.DataTypeString, entity.TypeOfCategory, filterOptions)
	if err != nil {
		return err
	}

	categoryUUID := r.URL.Query().Get(entity.CategoryUUID)
	filterOptions, err = processParam(categoryUUID, filter.DataTypeString, entity.CategoryUUID, filterOptions)
	if err != nil {
		return err
	}

	description := r.URL.Query().Get(entity.Description)
	filterOptions, err = processParam(description, filter.DataTypeString, entity.Description, filterOptions)
	if err != nil {
		return err
	}

	moneySum := r.URL.Query().Get(entity.MoneySum)
	filterOptions, err = processParam(moneySum, filter.DataTypeFloat, entity.MoneySum, filterOptions)
	if err != nil {
		return err
	}

	dateTime := r.URL.Query().Get(entity.DateTime)
	filterOptions, err = processParam(dateTime, filter.DataTypeDate, entity.DateTime, filterOptions)
	if err != nil {
		return err
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

func processParam(param, paramType, fieldName string, options filter.Options) (filter.Options, error) {
	validationErr := apperror.BadRequestError("filter params validation failed")
	if param != "" {
		operator := filter.OperatorEqual
		value := param

		if strings.Index(param, ":") != -1 {
			split := strings.Split(param, ":")
			operator = split[0]
			value = split[1]
		}

		values := strings.Split(value, ",")

		err := options.AddField(fieldName, operator, values, paramType)
		if err != nil {
			validationErr.WithParams(map[string]string{
				fieldName: err.Error(),
			})
			return options, validationErr
		}
	}
	return options, nil
}
