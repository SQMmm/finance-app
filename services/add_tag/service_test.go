package add_tag

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
)

func Test_service_AddTag(t *testing.T) {
	type data struct {
		addTagData
	}
	type mocks struct {
		addTagMock
	}
	type args struct {
		ctx context.Context
		tag *entities.Tag
	}
	tests := []struct {
		name    string
		args    args
		want    *entities.Tag
		wantErr bool
		data
		mocks
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				tag: &entities.Tag{
					ID:    134,
					User:  &entities.User{ID: 73},
					Title: "some tag",
				},
			},
			data: data{
				addTagData: addTagData{tag: entities.Tag{
					ID:    134,
					User:  &entities.User{ID: 73},
					Title: "some tag",
				}},
			},
			mocks: mocks{addTagMock: addTagMock{
				id: 157,
			}},
			want: &entities.Tag{
				ID:    157,
				User:  &entities.User{ID: 73},
				Title: "some tag",
			},
		},
		{
			name: "failed to add tag",
			args: args{
				ctx: context.Background(),
				tag: &entities.Tag{
					ID:    134,
					User:  &entities.User{ID: 73},
					Title: "some tag",
				},
			},
			data: data{
				addTagData: addTagData{tag: entities.Tag{
					ID:    134,
					User:  &entities.User{ID: 73},
					Title: "some tag",
				}},
			},
			mocks: mocks{addTagMock: addTagMock{
				err: errors.New("failed"),
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := &tagRepositoryMock{
				addTagMock: tt.mocks.addTagMock,
			}
			s := NewService(&logger.LoggerManagerMock{}, tags)
			got, err := s.AddTag(tt.args.ctx, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.AddTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.EqualValues(t, tt.want, got, "check result")
			assert.EqualValues(t, tt.data.addTagData, tags.addTagData, "check add tag data")
		})
	}
}
