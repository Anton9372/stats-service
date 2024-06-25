package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"stats-service/internal/apperror"
	"stats-service/internal/operations/entity"
	"stats-service/internal/operations/storage"
	"stats-service/pkg/logging"
	"stats-service/pkg/postgresql"
	"stats-service/pkg/utils"
	"time"
)

const queryWaitTime = 5 * time.Second

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) storage.Repository {
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

func (r *repository) FindAll(ctx context.Context, sortOptions storage.SortOptions) ([]entity.Operation, error) {
	qb := squirrel.Select("id, category_id, money_sum, description, date_time").From("public.operations")
	if sortOptions != nil {
		qb = qb.OrderBy(sortOptions.GetOrderBy())
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
