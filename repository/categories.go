package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"strings"
	"time"
)

type categories struct {
	wConn   *sql.DB
	rConn   *sql.DB
	timeout time.Duration
}

func NewCategories(w, r *sql.DB, t time.Duration) *categories {
	return &categories{
		wConn:   w,
		rConn:   r,
		timeout: t,
	}
}

func (c *categories) Add(ctx context.Context, category *entities.Category) (int64, error) {
	ctx, _ = context.WithTimeout(ctx, c.timeout)
	args := []interface{}{category.User.ID, category.Title}
	argStr := []string{"?", "?"}
	names := []string{"user_id", "title"}
	if category.Icon != nil {
		args = append(args, category.Icon.ID)
		argStr = append(argStr, "?")
		names = append(names, "icon_id")
	}
	if category.Color != nil {
		args = append(args, category.Color.ID)
		argStr = append(argStr, "?")
		names = append(names, "color_id")
	}

	res, err := c.wConn.ExecContext(ctx, fmt.Sprintf(`
insert into categories (%s)
values (%s)
`, strings.Join(names, ", "), strings.Join(argStr, ", ")), args...)
	if err != nil {
		return 0, fmt.Errorf("failed to do exec: %s", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %s", err)
	}
	return id, nil
}
