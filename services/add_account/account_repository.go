package add_account

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type AccountRepository interface {
	Add(ctx context.Context, account *entities.Account) (int64, error)
}
