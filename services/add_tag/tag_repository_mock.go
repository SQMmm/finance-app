package add_tag

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
)

type addTagData struct {
	tag entities.Tag
}
type addTagMock struct {
	id  int64
	err error
}
type tagRepositoryMock struct {
	addTagData

	addTagMock
}

func (trm *tagRepositoryMock) Add(ctx context.Context, tag *entities.Tag) (int64, error) {
	trm.addTagData = addTagData{tag: *tag}
	return trm.addTagMock.id, trm.addTagMock.err
}
