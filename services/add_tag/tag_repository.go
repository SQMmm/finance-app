package add_tag

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type tagRepository interface {
	Add(ctx context.Context, tag *entities.Tag) (int64, error)
}
