package add_account_group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
)

func Test_service_AddAccountGroup(t *testing.T) {
	type data struct {
		getAccountsData
		addData
	}
	type mocks struct {
		getAccountsMock
		addMock
	}
	type args struct {
		ctx   context.Context
		group *entities.AccountGroup
	}
	tests := []struct {
		name    string
		args    args
		want    *entities.AccountGroup
		wantErr bool
		data
		mocks
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:    124,
					User:  &entities.User{ID: 5656},
					Title: "new group",
					Accounts: []entities.Account{
						{ID: 63}, {ID: 899}, {ID: 244},
					},
				},
			},
			data: data{
				getAccountsData: getAccountsData{
					userID: 5656,
					ids:    []int64{63, 899, 244},
				},
				addData: addData{
					group: entities.AccountGroup{
						ID:    124,
						User:  &entities.User{ID: 5656},
						Title: "new group",
						Accounts: []entities.Account{
							{ID: 63, Title: "title1"}, {ID: 899, User: &entities.User{ID: 111}}, {ID: 244, StartBalance: 644, UseInReports: true},
						},
					},
				},
			},
			mocks: mocks{
				getAccountsMock: getAccountsMock{
					accounts: []entities.Account{
						{ID: 63, Title: "title1"}, {ID: 899, User: &entities.User{ID: 111}}, {ID: 244, StartBalance: 644, UseInReports: true},
					},
				},
				addMock: addMock{id: 147},
			},
			want: &entities.AccountGroup{
				ID:    147,
				User:  &entities.User{ID: 5656},
				Title: "new group",
				Accounts: []entities.Account{
					{ID: 63, Title: "title1"}, {ID: 899, User: &entities.User{ID: 111}}, {ID: 244, StartBalance: 644, UseInReports: true},
				},
			},
		},
		{
			name: "success there are not all accounts exist",
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:    124,
					User:  &entities.User{ID: 5656},
					Title: "new group",
					Accounts: []entities.Account{
						{ID: 63}, {ID: 899}, {ID: 244},
					},
				},
			},
			data: data{
				getAccountsData: getAccountsData{
					userID: 5656,
					ids:    []int64{63, 899, 244},
				},
				addData: addData{
					group: entities.AccountGroup{
						ID:    124,
						User:  &entities.User{ID: 5656},
						Title: "new group",
						Accounts: []entities.Account{
							{ID: 63}, {ID: 899},
						},
					},
				},
			},
			mocks: mocks{
				getAccountsMock: getAccountsMock{
					accounts: []entities.Account{
						{ID: 63}, {ID: 899},
					},
				},
				addMock: addMock{id: 368},
			},
			want: &entities.AccountGroup{
				ID:    368,
				User:  &entities.User{ID: 5656},
				Title: "new group",
				Accounts: []entities.Account{
					{ID: 63}, {ID: 899},
				},
			},
		},
		{
			name: "success without accounts",
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:       124,
					User:     &entities.User{ID: 5656},
					Title:    "new group",
					Accounts: []entities.Account{},
				},
			},
			data: data{
				addData: addData{
					group: entities.AccountGroup{
						ID:       124,
						User:     &entities.User{ID: 5656},
						Title:    "new group",
						Accounts: []entities.Account{},
					},
				},
			},
			mocks: mocks{
				addMock: addMock{id: 368},
			},
			want: &entities.AccountGroup{
				ID:       368,
				User:     &entities.User{ID: 5656},
				Title:    "new group",
				Accounts: []entities.Account{},
			},
		},
		{
			name: "failed to get accounts",
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:    124,
					User:  &entities.User{ID: 5656},
					Title: "new group",
					Accounts: []entities.Account{
						{ID: 63}, {ID: 899}, {ID: 244},
					},
				},
			},
			data: data{
				getAccountsData: getAccountsData{
					userID: 5656,
					ids:    []int64{63, 899, 244},
				},
			},
			mocks: mocks{
				getAccountsMock: getAccountsMock{
					err: errors.New("failed"),
				},
			},
			wantErr: true,
		},
		{
			name: "failed to add group",
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:    124,
					User:  &entities.User{ID: 5656},
					Title: "new group",
					Accounts: []entities.Account{
						{ID: 63}, {ID: 899}, {ID: 244},
					},
				},
			},
			data: data{
				getAccountsData: getAccountsData{
					userID: 5656,
					ids:    []int64{63, 899, 244},
				},
				addData: addData{
					group: entities.AccountGroup{
						ID:    124,
						User:  &entities.User{ID: 5656},
						Title: "new group",
						Accounts: []entities.Account{
							{ID: 63}, {ID: 899}, {ID: 244},
						},
					},
				},
			},
			mocks: mocks{
				getAccountsMock: getAccountsMock{
					accounts: []entities.Account{
						{ID: 63}, {ID: 899}, {ID: 244},
					},
				},
				addMock: addMock{err: errors.New("failed")},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &accountsRepositoryMock{
				getAccountsMock: tt.mocks.getAccountsMock,
			}
			group := &groupRepositoryMock{
				addMock: tt.mocks.addMock,
			}
			s := NewService(&logger.LoggerManagerMock{}, account, group)
			got, err := s.AddAccountGroup(tt.args.ctx, tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.AddAccountGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, tt.want, got, "check result")
			assert.EqualValues(t, tt.data.getAccountsData, account.getAccountsData, "check get accounts data")
			assert.EqualValues(t, tt.data.addData, group.addData, "check add group data")
		})
	}
}
