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
		ID:           a.ID,
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

type AccountGroup struct {
	ID       int64     `json:"id"`
	UserID   int64     `json:"user_id"`
	Title    string    `json:"title"`
	Accounts []Account `json:"accounts,omitempty"`
}

func (ag *AccountGroup) Validate() error {
	err := validation.Errors{
		"title length is incorrect": validation.Validate(
			ag.Title,
			validation.Required,
			validation.RuneLength(1, 120),
		),
	}
	return err.Filter()
}

func (ag *AccountGroup) GetEntity() *entities.AccountGroup {
	g := &entities.AccountGroup{
		Title:    ag.Title,
		Accounts: make([]entities.Account, len(ag.Accounts)),
	}
	for i, acc := range ag.Accounts {
		g.Accounts[i] = *acc.GetEntity()
	}

	return g
}

func GetAccountGroupFromEntity(g *entities.AccountGroup) *AccountGroup {
	ag := &AccountGroup{
		ID:       g.ID,
		UserID:   g.User.ID,
		Title:    g.Title,
		Accounts: make([]Account, len(g.Accounts)),
	}
	for i, acc := range g.Accounts {
		ag.Accounts[i] = *GetAccountFromEntity(&acc)
	}

	return ag
}
