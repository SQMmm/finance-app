package add_account

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type addAccountData struct {
	account entities.Account
}
type addAccountMock struct {
	id  int64
	err error
}
type accountRepositoryMock struct {
	addAccountData

	addAccountMock
}

func (arm *accountRepositoryMock) Add(_ context.Context, acc *entities.Account) (int64, error) {
	arm.addAccountData = addAccountData{account: *acc}

	return arm.addAccountMock.id, arm.addAccountMock.err
}
