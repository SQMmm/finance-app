package add_category

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type getIconByIDData struct {
	id int64
}
type getColorByIDData struct {
	id int64
}
type getIconByIDMock struct {
	icon *entities.Icon
	err  error
}
type getColorByIDMock struct {
	color *entities.IconColor
	err   error
}
type iconRepositoryMock struct {
	getIconByIDData
	getColorByIDData

	getIconByIDMock
	getColorByIDMock
}

func (irm *iconRepositoryMock) GetIconByID(ctx context.Context, id int64) (*entities.Icon, error) {
	irm.getIconByIDData = getIconByIDData{id: id}
	return irm.getIconByIDMock.icon, irm.getIconByIDMock.err
}

func (irm *iconRepositoryMock) GetColorByID(ctx context.Context, id int64) (*entities.IconColor, error) {
	irm.getColorByIDData = getColorByIDData{id: id}
	return irm.getColorByIDMock.color, irm.getColorByIDMock.err
}
