package add_account

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
)

func Test_service_AddAccount(t *testing.T) {
	type data struct {
		addAccountData
	}
	type mocks struct {
		addAccountMock
	}
	type args struct {
		ctx     context.Context
		account *entities.Account
	}
	tests := []struct {
		name    string
		args    args
		want    *entities.Account
		wantErr bool
		data
		mocks
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				account: &entities.Account{
					ID:           1243,
					User:         &entities.User{ID: 123},
					Title:        "vfjvf",
					StartBalance: 98.1,
					UseInReports: true,
				},
			},
			data: data{addAccountData: addAccountData{account: entities.Account{
				ID:           1243,
				User:         &entities.User{ID: 123},
				Title:        "vfjvf",
				StartBalance: 98.1,
				UseInReports: true,
			}}},
			mocks: mocks{addAccountMock: addAccountMock{id: 467}},
			want: &entities.Account{
				ID:           467,
				User:         &entities.User{ID: 123},
				Title:        "vfjvf",
				StartBalance: 98.1,
				UseInReports: true,
			},
		},
		{
			name: "failed to add",
			args: args{
				ctx: context.Background(),
				account: &entities.Account{
					ID:           1243,
					User:         &entities.User{ID: 123},
					Title:        "vfjvf",
					StartBalance: 98.1,
					UseInReports: true,
				},
			},
			data: data{addAccountData: addAccountData{account: entities.Account{
				ID:           1243,
				User:         &entities.User{ID: 123},
				Title:        "vfjvf",
				StartBalance: 98.1,
				UseInReports: true,
			}}},
			mocks:   mocks{addAccountMock: addAccountMock{err: errors.New("failed")}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accounts := &accountRepositoryMock{
				addAccountMock: tt.mocks.addAccountMock,
			}
			s := NewService(&logger.LoggerManagerMock{}, accounts)
			got, err := s.AddAccount(tt.args.ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.AddAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.EqualValues(t, tt.want, got, "check result")
			assert.EqualValues(t, tt.data.addAccountData, accounts.addAccountData, "check add account data")
		})
	}
}
