package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
	"strings"
	"time"
)

type accountGroups struct {
	manager logger.LoggerManager
	wConn   *sql.DB
	rConn   *sql.DB
	timeout time.Duration
}

func NewAccountGroups(w, r *sql.DB, t time.Duration, m logger.LoggerManager) *accountGroups {
	return &accountGroups{
		wConn:   w,
		rConn:   r,
		timeout: t,
		manager: m,
	}
}

func (ag *accountGroups) Add(ctx context.Context, group *entities.AccountGroup) (id int64, err error) {
	ctx, _ = context.WithTimeout(ctx, ag.timeout)

	tx, err := ag.wConn.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin context: %s", err)
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
			return
		}
		if newErr := tx.Rollback(); newErr != nil {
			ag.manager.LogCtx(ctx).Errorf("failed to rollback transaction: %s", err)
		}
	}()

	result, err := tx.ExecContext(ctx, `
insert into account_groups (user_id, title) 
values (?, ?)
`, group.User.ID, group.Title)
	if err != nil {
		return 0, fmt.Errorf("failed to add group: %s", err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %s", err)
	}

	if len(group.Accounts) == 0 {
		return id, nil
	}

	ctx = context.Background()
	ctx, _ = context.WithTimeout(ctx, ag.timeout)

	values := make([]string, len(group.Accounts))
	args := make([]interface{}, len(group.Accounts)*3)
	for i, acc := range group.Accounts {
		values[i] = "(?, ?, ?)"
		args[i*3] = group.User.ID
		args[i*3+1] = acc.ID
		args[i*3+2] = id
	}

	_, err = tx.ExecContext(ctx,
		fmt.Sprintf("insert into `groups` (user_id, account_id, group_id) values %s", strings.Join(values, ", ")),
		args...)
	if err != nil {
		return 0, fmt.Errorf("failed to add links: %s", err)
	}

	return id, nil
}
