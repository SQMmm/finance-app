package add_account_group

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type addData struct {
	group entities.AccountGroup
}
type addMock struct {
	id  int64
	err error
}
type groupRepositoryMock struct {
	addData

	addMock
}

func (grm *groupRepositoryMock) Add(ctx context.Context, group *entities.AccountGroup) (int64, error) {
	grm.addData = addData{group: *group}
	return grm.addMock.id, grm.addMock.err
}
