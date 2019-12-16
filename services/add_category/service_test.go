package add_category

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
)

func Test_service_Add(t *testing.T) {
	type data struct {
		addData
		getIconByIDData
		getColorByIDData
	}
	type mocks struct {
		addMock
		getIconByIDMock
		getColorByIDMock
	}
	type args struct {
		ctx      context.Context
		category *entities.Category
	}
	tests := []struct {
		name    string
		args    args
		want    *entities.Category
		wantErr bool
		data
		mocks
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				category: &entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
					Icon:  &entities.Icon{ID: 124},
					Color: &entities.IconColor{ID: 85475},
				},
			},
			data: data{
				getIconByIDData:  getIconByIDData{id: 124},
				getColorByIDData: getColorByIDData{id: 85475},
				addData: addData{category: entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
					Icon:  &entities.Icon{ID: 124},
					Color: &entities.IconColor{ID: 85475},
				}},
			},
			mocks: mocks{
				addMock:          addMock{id: 98},
				getIconByIDMock:  getIconByIDMock{icon: &entities.Icon{ID: 124, Path: "path", Name: "name icon"}},
				getColorByIDMock: getColorByIDMock{color: &entities.IconColor{ID: 85475, Color: "some color", Name: "color name"}},
			},
			want: &entities.Category{
				ID:    98,
				User:  &entities.User{ID: 53},
				Title: "category title",
				Icon:  &entities.Icon{ID: 124, Path: "path", Name: "name icon"},
				Color: &entities.IconColor{ID: 85475, Color: "some color", Name: "color name"},
			},
		},
		{
			name: "success without icon and color",
			args: args{
				ctx: context.Background(),
				category: &entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
				},
			},
			data: data{
				addData: addData{category: entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
				}},
			},
			mocks: mocks{
				addMock: addMock{id: 98},
			},
			want: &entities.Category{
				ID:    98,
				User:  &entities.User{ID: 53},
				Title: "category title",
			},
		},
		{
			name: "failed to add category",
			args: args{
				ctx: context.Background(),
				category: &entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
					Icon:  &entities.Icon{ID: 124},
					Color: &entities.IconColor{ID: 85475},
				},
			},
			data: data{
				addData: addData{category: entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
					Icon:  &entities.Icon{ID: 124},
					Color: &entities.IconColor{ID: 85475},
				}},
			},
			mocks: mocks{
				addMock: addMock{err: errors.New("failed")},
			},
			wantErr: true,
		},
		{
			name: "failed to get icon",
			args: args{
				ctx: context.Background(),
				category: &entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
					Icon:  &entities.Icon{ID: 124},
					Color: &entities.IconColor{ID: 85475},
				},
			},
			data: data{
				getIconByIDData: getIconByIDData{id: 124},
				addData: addData{category: entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
					Icon:  &entities.Icon{ID: 124},
					Color: &entities.IconColor{ID: 85475},
				}},
			},
			mocks: mocks{
				addMock:         addMock{id: 98},
				getIconByIDMock: getIconByIDMock{err: errors.New("failed")},
			},
			wantErr: true,
		},
		{
			name: "failed to get color",
			args: args{
				ctx: context.Background(),
				category: &entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
					Icon:  &entities.Icon{ID: 124},
					Color: &entities.IconColor{ID: 85475},
				},
			},
			data: data{
				getIconByIDData:  getIconByIDData{id: 124},
				getColorByIDData: getColorByIDData{id: 85475},
				addData: addData{category: entities.Category{
					ID:    124,
					User:  &entities.User{ID: 53},
					Title: "category title",
					Icon:  &entities.Icon{ID: 124},
					Color: &entities.IconColor{ID: 85475},
				}},
			},
			mocks: mocks{
				addMock:          addMock{id: 98},
				getIconByIDMock:  getIconByIDMock{icon: &entities.Icon{ID: 124, Path: "path", Name: "name icon"}},
				getColorByIDMock: getColorByIDMock{err: errors.New("failed")},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := &categoryRepositoryMock{
				addMock: tt.mocks.addMock,
			}
			icon := &iconRepositoryMock{
				getIconByIDMock:  tt.mocks.getIconByIDMock,
				getColorByIDMock: tt.mocks.getColorByIDMock,
			}
			s := NewService(&logger.LoggerManagerMock{}, category, icon)
			got, err := s.AddCategory(tt.args.ctx, tt.args.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, tt.want, got, "check result")
			assert.EqualValues(t, tt.data.getIconByIDData, icon.getIconByIDData, "check get icon data")
			assert.EqualValues(t, tt.data.getColorByIDData, icon.getColorByIDData, "check get color data")
			assert.EqualValues(t, tt.data.addData, category.addData, "check add category data")
		})
	}
}
