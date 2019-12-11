package data

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/sqmmm/finance-app/entities"
)

type Account struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	UserID       int64   `json:"user_id"`
	StartBalance float64 `json:"start_balance"`
	UseInReports bool    `json:"use_in_reports"`
}

func (a *Account) Validate() error {
	err := validation.Errors{
		"title length is incorrect": validation.Validate(
			a.Title,
			validation.Required,
			validation.RuneLength(1, 120),
		),
	}
	return err.Filter()
}

func (a *Account) GetEntity() *entities.Account {
	return &entities.Account{
		Title:        a.Title,
		StartBalance: a.StartBalance,
		UseInReports: a.UseInReports,
	}
}

func GetAccountFromEntity(acc *entities.Account) *Account {
	return &Account{
		ID:           acc.ID,
		Title:        acc.Title,
		UserID:       acc.User.ID,
		StartBalance: acc.StartBalance,
		UseInReports: acc.UseInReports,
	}
}
