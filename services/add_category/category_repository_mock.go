package add_category

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type addData struct {
	category entities.Category
}
type addMock struct {
	id  int64
	err error
}
type categoryRepositoryMock struct {
	addData

	addMock
}

func (crm *categoryRepositoryMock) Add(ctx context.Context, category *entities.Category) (int64, error) {
	crm.addData = addData{category: *category}
	return crm.addMock.id, crm.addMock.err
}
