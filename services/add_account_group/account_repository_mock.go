package add_account_group

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type getAccountsData struct {
	userID int64
	ids    []int64
}
type getAccountsMock struct {
	accounts []entities.Account
	err      error
}
type accountsRepositoryMock struct {
	getAccountsData

	getAccountsMock
}

func (arm *accountsRepositoryMock) GetUserAccountsByIDs(ctx context.Context, userID int64, ids []int64) (entities.Accounts, error) {
	arm.getAccountsData = getAccountsData{userID: userID, ids: ids}
	return arm.getAccountsMock.accounts, arm.getAccountsMock.err
}
