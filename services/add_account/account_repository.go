package add_account

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type accountRepository interface {
	Add(ctx context.Context, account *entities.Account) (int64, error)
}
