package add_category

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type iconRepository interface {
	GetIconByID(ctx context.Context, id int64) (*entities.Icon, error)
	GetColorByID(ctx context.Context, id int64) (*entities.IconColor, error)
}

