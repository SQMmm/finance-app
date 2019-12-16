package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"time"
)

type icons struct {
	wConn   *sql.DB
	rConn   *sql.DB
	timeout time.Duration
}

func NewIcons(w, r *sql.DB, t time.Duration) *icons {
	return &icons{
		wConn:   w,
		rConn:   r,
		timeout: t,
	}
}

func (i *icons) GetIconByID(ctx context.Context, id int64) (*entities.Icon, error) {
	icon := &entities.Icon{}
	ctx, _ = context.WithTimeout(ctx, i.timeout)

	row := i.rConn.QueryRowContext(ctx, `
select id, name, path
from icons
where id = ?`, id)

	if err := row.Scan(&icon.ID, &icon.Name, &icon.Path); err != nil {
		return nil, fmt.Errorf("faield to do query: %s", err)
	}

	return icon, nil
}

func (i *icons) GetColorByID(ctx context.Context, id int64) (*entities.IconColor, error) {
	color := &entities.IconColor{}
	ctx, _ = context.WithTimeout(ctx, i.timeout)

	row := i.rConn.QueryRowContext(ctx, `
select id, name, color
from icon_colors
where id = ?`, id)

	if err := row.Scan(&color.ID, &color.Name, &color.Color); err != nil {
		return nil, fmt.Errorf("faield to do query: %s", err)
	}

	return color, nil
}
