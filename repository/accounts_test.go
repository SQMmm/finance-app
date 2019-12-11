package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
	"time"

	"github.com/sqmmm/finance-app/entities"
)

var ctxCanceled context.Context

func init() {
	var cancel context.CancelFunc
	ctxCanceled, cancel = context.WithCancel(context.Background())
	cancel()
}

func Test_accounts_Add(t *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	type args struct {
		ctx     context.Context
		account *entities.Account
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
		delay   time.Duration
	}{
		{
			name:   "success",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				account: &entities.Account{
					User:         &entities.User{ID: 33},
					Title:        "title",
					StartBalance: 869.76,
					UseInReports: true,
				},
			},
			want:  74,
			delay: 10 * time.Millisecond,
		},
		{
			name:   "failed by timeout",
			fields: fields{timeout: time.Millisecond},
			args: args{
				ctx: context.Background(),
				account: &entities.Account{
					User:         &entities.User{ID: 33},
					Title:        "title",
					StartBalance: 869.76,
					UseInReports: true,
				},
			},
			wantErr: true,
			delay:   10 * time.Millisecond,
		},
		{
			name:   "failed by canceled context",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: ctxCanceled,
				account: &entities.Account{
					User:         &entities.User{ID: 33},
					Title:        "title",
					StartBalance: 869.76,
					UseInReports: true,
				},
			},
			wantErr: true,
			delay:   10 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalln(err)
			}
			acc := tt.args.account

			mock.ExpectExec("insert into accounts").
				WithArgs(acc.User.ID, acc.Title, acc.StartBalance, acc.UseInReports).
				WillReturnResult(sqlmock.NewResult(tt.want, 1)).
				WillDelayFor(tt.delay)

			a := NewAccounts(db, &sql.DB{}, tt.fields.timeout)
			got, err := a.Add(tt.args.ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("accounts.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("accounts.Add() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); !tt.wantErr && err != nil {
				t.Errorf("accounts.Add() there were unfulfilled expectations: %s", err)
			}
		})
	}
}
