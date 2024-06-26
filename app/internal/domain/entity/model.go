package entity

import "time"

type CategoryType string

const (
	IncomeType  CategoryType = "Income"
	ExpenseType CategoryType = "Expense"
)

type Category struct {
	UUID     string       `json:"uuid"`
	UserUUID string       `json:"user_uuid"`
	Name     string       `json:"name"`
	Type     CategoryType `json:"type"`
}

type Operation struct {
	UUID         string    `json:"uuid"`
	CategoryUUID string    `json:"category_uuid"`
	Description  string    `json:"description"`
	MoneySum     float64   `json:"money_sum"`
	DateTime     time.Time `json:"date_time"`
}

// columns in public.operations
const (
	UUID         = "id"
	CategoryUUID = "category_id"
	Description  = "description"
	MoneySum     = "money_sum"
	DateTime     = "date_time"
)

type Report struct {
	TotalMoneySum float64     `json:"total_money_sum"`
	Operations    []Operation `json:"operations"`
}

func NewReport(operations []Operation) Report {
	sum := 0.0
	for _, op := range operations {
		sum += op.MoneySum
	}
	return Report{
		TotalMoneySum: sum,
		Operations:    operations,
	}
}
