package add_category

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type categoryRepository interface {
	Add(ctx context.Context, category *entities.Category) (int64, error)
}
