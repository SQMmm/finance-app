package add_account_group

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type groupRepository interface {
	Add(ctx context.Context, group *entities.AccountGroup) (int64, error)
}
