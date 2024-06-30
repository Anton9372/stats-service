package filter

import (
	"fmt"
	"strconv"
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
	Values   []string
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

func (o *options) AddField(name, operator string, values []string, dataType string) error {
	field := Field{
		Name:     name,
		Values:   values,
		Operator: operator,
		DataType: dataType,
	}

	if err := validateField(field); err != nil {
		return err
	}

	o.fields = append(o.fields, field)
	return nil
}

func validateField(field Field) error {
	if err := validateOperator(field.Operator); err != nil {
		return err
	}

	for _, value := range field.Values {
		if err := validateValue(value, field.DataType); err != nil {
			return err
		}
	}

	if err := validateOperatorAndDataType(field.Operator, field.DataType); err != nil {
		return err
	}

	if err := validateOperatorAndValues(field.Operator, field.Values); err != nil {
		return err
	}

	return nil
}

func validateOperatorAndValues(operator string, values []string) error {
	switch operator {
	case OperatorEqual:
	case OperatorNotEqual:
	case OperatorLowerThan:
		if len(values) > 1 {
			return fmt.Errorf("with the '%s' operator only one value can be used", OperatorLowerThan)
		}
	case OperatorLowerThanEqual:
		if len(values) > 1 {
			return fmt.Errorf("with the '%s' operator only one value can be used", OperatorLowerThanEqual)
		}
	case OperatorGreaterThan:
		if len(values) > 1 {
			return fmt.Errorf("with the '%s' operator only one value can be used", OperatorGreaterThan)
		}
	case OperatorGreaterThanEqual:
		if len(values) > 1 {
			return fmt.Errorf("with the '%s' operator only one value can be used", OperatorGreaterThanEqual)
		}
	case OperatorBetween:
		if len(values) != 2 {
			return fmt.Errorf("with the '%s' operator two values should be used", OperatorBetween)
		}
	case OperatorSubString:
		if len(values) > 1 {
			return fmt.Errorf("with the '%s' operator only one value can be used", OperatorSubString)
		}
	}
	return nil
}

func validateOperatorAndDataType(operator, dataType string) error {
	switch dataType {
	case DataTypeString:
		if operator != OperatorEqual && operator != OperatorSubString {
			return fmt.Errorf("with the string data type only the '=' and '%s' operator can be used", OperatorSubString)
		}
	case DataTypeFloat:
		if operator == OperatorSubString {
			return fmt.Errorf("with the float data type the '%s' operator can not be used", OperatorSubString)
		}
	case DataTypeDate:
		if operator != OperatorEqual && operator != OperatorBetween {
			return fmt.Errorf("with the date data type only the '=' and '%s' operator can be used", OperatorBetween)
		}
	}
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
		_, err := time.Parse(time.DateOnly, value)
		if err != nil {
			return fmt.Errorf("date should be in format %s", time.DateOnly)
		}
		return nil
	default:
		return fmt.Errorf("invalid data type: %s", dataType)
	}
}
