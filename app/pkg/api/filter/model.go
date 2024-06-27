package filter

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DataTypeString = "string"
	DataTypeFloat  = "float"
	DataTypeDate   = "date"

	OperatorEqual            = "eq"
	OperatorNotEqual         = "neq"
	OperatorLowerThan        = "lt"
	OperatorLowerThanEqual   = "lte"
	OperatorGreaterThan      = "gt"
	OperatorGreaterThanEqual = "gte"
	OperatorBetween          = "between"
	OperatorSubString        = "substr"
)

type options struct {
	limit  int
	fields []Field
}

type Field struct {
	Name     string
	Operator string
	Value    string
	DataType string
}

func NewOptions(limit int) Options {
	return &options{
		limit: limit,
	}
}

func (o *options) Fields() []Field {
	return o.fields
}

func (o *options) Limit() int {
	return o.limit
}

func (o *options) AddField(name, operator, value, dataType string) error {
	if err := validateOperator(operator); err != nil {
		return err
	}
	if err := validateValue(value, dataType); err != nil {
		return err
	}
	o.fields = append(o.fields, Field{
		Name:     name,
		Value:    value,
		Operator: operator,
		DataType: dataType,
	})
	return nil
}

func validateOperator(operator string) error {
	switch operator {
	case OperatorEqual:
	case OperatorNotEqual:
	case OperatorLowerThan:
	case OperatorLowerThanEqual:
	case OperatorGreaterThan:
	case OperatorGreaterThanEqual:
	case OperatorBetween:
	case OperatorSubString:
	default:
		return fmt.Errorf("invalid operator: %s", operator)
	}
	return nil
}

func validateValue(value, dataType string) error {
	switch dataType {
	case DataTypeString:
		return nil
	case DataTypeFloat:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("failed converting value to number")
		}
		return nil
	case DataTypeDate:
		split := strings.Split(value, ":")
		_, err0 := time.Parse(time.DateOnly, split[0])
		_, err1 := time.Parse(time.DateOnly, split[1])
		if err0 != nil || err1 != nil {
			return fmt.Errorf("date should be in format %s or %s:%s", time.DateOnly, time.DateOnly, time.DateOnly)
		}
		return nil
	default:
		return fmt.Errorf("invalid data type: %s", dataType)
	}
}
