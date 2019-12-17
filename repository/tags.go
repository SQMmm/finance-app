package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"time"
)

type tags struct {
	wConn   *sql.DB
	rConn   *sql.DB
	timeout time.Duration
}

func NewTags(w, r *sql.DB, t time.Duration) *tags {
	return &tags{
		wConn:   w,
		rConn:   r,
		timeout: t,
	}
}

func (t *tags) Add(ctx context.Context, tag *entities.Tag) (int64, error) {
	ctx, _ = context.WithTimeout(ctx, t.timeout)

	res, err := t.wConn.ExecContext(ctx, `
insert into tags(user_id, title)
values (?, ?)
`, tag.User.ID, tag.Title)
	if err != nil {
		return 0, fmt.Errorf("faile to do exec: %s", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %s", err)
	}

	return id, nil
}
