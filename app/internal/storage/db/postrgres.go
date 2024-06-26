package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"stats-service/internal/apperror"
	"stats-service/internal/domain/entity"
	"stats-service/internal/domain/service"
	"stats-service/internal/storage/sorting"
	"stats-service/pkg/api/filter"
	"stats-service/pkg/logging"
	"stats-service/pkg/postgresql"
	"stats-service/pkg/utils"
	"strings"
	"time"
)

const queryWaitTime = 5 * time.Second

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) service.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func handleSQLError(err error, logger *logging.Logger) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
			pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
		logger.Error(newErr)

		if pgErr.Code == "23505" { //uniqueness violation
			return apperror.BadRequestError("")
		}
		return newErr
	}

	return err
}

func processFilterOptionsWithSquirrel(qb squirrel.SelectBuilder, options filter.Options) (squirrel.SelectBuilder, error) {
	fields := options.Fields()

	for _, field := range fields {
		switch field.Operator {
		case filter.OperatorEqual:
			qb = qb.Where(squirrel.Eq{field.Name: field.Value})
		case filter.OperatorNotEqual:
			qb = qb.Where(squirrel.NotEq{field.Name: field.Value})
		case filter.OperatorLowerThan:
			qb = qb.Where(squirrel.Lt{field.Name: field.Value})
		case filter.OperatorLowerThanEqual:
			qb = qb.Where(squirrel.LtOrEq{field.Name: field.Value})
		case filter.OperatorGreaterThan:
			qb = qb.Where(squirrel.Gt{field.Name: field.Value})
		case filter.OperatorGreaterThanEqual:
			qb = qb.Where(squirrel.GtOrEq{field.Name: field.Value})
		case filter.OperatorBetween:
			values := strings.Split(field.Value, ":")
			if len(values) != 2 {
				return qb, fmt.Errorf("between operator requires two values")
			}
			qb = qb.Where(squirrel.Expr(fmt.Sprintf("%s BETWEEN ? AND ?", field.Name),
				fmt.Sprintf("%s 00:00:00", values[0]), fmt.Sprintf("%s 23:59:59", values[1])))
		case filter.OperatorSubString:
			qb = qb.Where(squirrel.Like{field.Name: "%" + field.Value + "%"})
		default:
			return qb, fmt.Errorf("invalid operator: %s", field.Operator)
		}
	}

	qb = qb.PlaceholderFormat(squirrel.Dollar)
	return qb, nil
}

func (r *repository) FindAll(ctx context.Context, sortOptions sorting.SortOptions, filterOptions filter.Options) ([]entity.Operation, error) {
	var err error
	qb := squirrel.Select("id, category_id, money_sum, description, date_time").From("public.operations")

	if sortOptions != nil {
		qb = qb.OrderBy(sortOptions.GetOrderBy())
	}

	if filterOptions != nil {
		qb, err = processFilterOptionsWithSquirrel(qb, filterOptions)
		if err != nil {
			return nil, err
		}
	}

	sql, i, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query into a SQL string: %w", err)
	}
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", utils.FormatSQLQuery(sql)))

	nCtx, cancel := context.WithTimeout(ctx, queryWaitTime)
	defer cancel()
	rows, err := r.client.Query(nCtx, sql, i...)
	if err != nil {
		return nil, handleSQLError(err, r.logger)
	}
	defer rows.Close()

	operations := make([]entity.Operation, 0)
	for rows.Next() {
		var op entity.Operation
		err = rows.Scan(&op.UUID, &op.CategoryUUID, &op.MoneySum, &op.Description, &op.DateTime)
		if err != nil {
			return nil, err
		}
		operations = append(operations, op)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return operations, nil
}