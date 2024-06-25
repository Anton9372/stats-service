package entity

import "time"

//By category type
//By category
//By date

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
