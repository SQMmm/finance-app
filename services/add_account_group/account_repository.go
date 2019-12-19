package add_account_group

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type accountRepository interface {
	GetUserAccountsByIDs(ctx context.Context, userID int64, ids []int64) (entities.Accounts, error)
}
