package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"time"
)

type accounts struct {
	wConn   *sql.DB
	rConn   *sql.DB
	timeout time.Duration
}

func NewAccounts(w, r *sql.DB, t time.Duration) *accounts {
	return &accounts{
		wConn:   w,
		rConn:   r,
		timeout: t,
	}
}

func (a *accounts) Add(ctx context.Context, account *entities.Account) (int64, error) {
	ctx, _ = context.WithTimeout(ctx, a.timeout)

	result, err := a.wConn.ExecContext(ctx, `
insert into accounts (user_id, title, start_balance, use_in_report)
values (?, ?, ?, ?)
`, account.User.ID, account.Title, account.StartBalance, account.UseInReports)
	if err != nil {
		return 0, fmt.Errorf("failed to do exec: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %s", err)
	}

	return id, nil
}
