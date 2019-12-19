package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"strings"
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

func (a *accounts) GetUserAccountsByIDs(ctx context.Context, userID int64, ids []int64) (entities.Accounts, error) {
	ctx, _ = context.WithTimeout(ctx, a.timeout)

	accounts := make([]entities.Account, 0)

	args := make([]interface{}, len(ids)+1)
	strArgs := make([]string, len(ids))
	for i, id := range ids {
		args[i] = id
		strArgs[i] = "?"
	}
	args[len(args)-1] = userID
	rows, err := a.rConn.QueryContext(ctx, fmt.Sprintf(`
select id, title, start_balance, use_in_report from accounts 
where user_id = ? and id in (%s) 
`, strings.Join(strArgs, ",")), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to de query: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		acc := entities.Account{User: &entities.User{ID: userID}}
		if err := rows.Scan(&acc.ID, &acc.Title, &acc.StartBalance, &acc.UseInReports); err != nil {
			return nil, fmt.Errorf("failed to scan: %s", err)
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}
