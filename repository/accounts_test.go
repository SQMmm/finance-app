package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

func Test_accounts_GetUserAccountsByIDs(t *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	type args struct {
		ctx    context.Context
		userID int64
		ids    []int64
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     entities.Accounts
		wantErr  bool
		delay    time.Duration
		rowsMock [][]driver.Value
	}{
		{
			name:   "success",
			fields: fields{timeout: time.Second},
			args: args{
				ctx:    context.Background(),
				userID: 1245,
				ids:    []int64{1, 2, 3},
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{1, "title", 12, true},
				{2, "title", 12, false},
			},
			want: entities.Accounts{
				{ID: 1, Title: "title", User: &entities.User{ID: 1245}, StartBalance: 12, UseInReports: true},
				{ID: 2, Title: "title", User: &entities.User{ID: 1245}, StartBalance: 12},
			},
		},
		{
			name:   "failed by timeout",
			fields: fields{timeout: time.Millisecond},
			args: args{
				ctx:    context.Background(),
				userID: 1245,
				ids:    []int64{1, 2, 3},
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{1, "title", 12, true},
			},
			wantErr: true,
		},
		{
			name:   "failed by canceled context",
			fields: fields{timeout: time.Second},
			args: args{
				ctx:    ctxCanceled,
				userID: 1245,
				ids:    []int64{1, 2, 3},
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{1, "title", 12, true},
			},
			wantErr: true,
		},
		{
			name:   "failed by scan",
			fields: fields{timeout: time.Second},
			args: args{
				ctx:    context.Background(),
				userID: 1245,
				ids:    []int64{1, 2, 3},
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{1, "title", nil, true},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalln(err)
			}
			args := make([]driver.Value, len(tt.args.ids)+1)
			for i, id := range tt.args.ids {
				args[i] = id
			}
			args[len(args)-1] = tt.args.userID

			rows := sqlmock.NewRows([]string{"id", "title", "start_balance", "user_in_report"})
			for _, m := range tt.rowsMock {
				rows = rows.AddRow(m...)
			}

			mock.ExpectQuery(
				`select id, title, start_balance, use_in_report from accounts where user_id = \? and id in`).
				WithArgs(args...).
				RowsWillBeClosed().
				WillDelayFor(tt.delay).
				WillReturnRows(rows)

			a := &accounts{
				wConn:   &sql.DB{},
				rConn:   db,
				timeout: tt.fields.timeout,
			}
			got, err := a.GetUserAccountsByIDs(tt.args.ctx, tt.args.userID, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("accounts.GetUserAccountsByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.EqualValues(t, tt.want, got, "check result")

			if err := mock.ExpectationsWereMet(); !tt.wantErr && err != nil {
				t.Errorf("accounts.GetUserAccountsByIDs() there were unfulfilled expectations: %s", err)
			}
		})
	}
}
