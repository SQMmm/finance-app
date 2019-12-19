package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sqmmm/finance-app/internal/logger"
	"log"
	"testing"
	"time"

	"github.com/sqmmm/finance-app/entities"
)

func Test_accountGroups_Add(t *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	type args struct {
		ctx   context.Context
		group *entities.AccountGroup
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        int64
		wantErr     bool
		delay       time.Duration
		secondDelay time.Duration
		resultError bool
	}{
		{
			name:   "success",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:    274,
					User:  &entities.User{ID: 235},
					Title: "some title",
					Accounts: entities.Accounts{
						{ID: 135}, {ID: 55},
					},
				},
			},
			delay:       10 * time.Millisecond,
			secondDelay: 10 * time.Millisecond,
			want:        135,
		},
		{
			name:   "success without accounts",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:    274,
					User:  &entities.User{ID: 235},
					Title: "some title",
				},
			},
			delay:       10 * time.Millisecond,
			secondDelay: 10 * time.Millisecond,
			want:        135,
		},
		{
			name:   "failed by timeout",
			fields: fields{timeout: time.Millisecond},
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:    274,
					User:  &entities.User{ID: 235},
					Title: "some title",
					Accounts: entities.Accounts{
						{ID: 135}, {ID: 55},
					},
				},
			},
			delay:   10 * time.Millisecond,
			wantErr: true,
		},
		{
			name:   "failed by canceled timeout",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: ctxCanceled,
				group: &entities.AccountGroup{
					ID:    274,
					User:  &entities.User{ID: 235},
					Title: "some title",
					Accounts: entities.Accounts{
						{ID: 135}, {ID: 55},
					},
				},
			},
			delay:   10 * time.Millisecond,
			wantErr: true,
		},
		{
			name:   "failed by getting last insert id",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:    274,
					User:  &entities.User{ID: 235},
					Title: "some title",
					Accounts: entities.Accounts{
						{ID: 135}, {ID: 55},
					},
				},
			},
			delay:       10 * time.Millisecond,
			wantErr:     true,
			resultError: true,
		},
		{
			name:   "failed by second timeout",
			fields: fields{timeout: 50 * time.Millisecond},
			args: args{
				ctx: context.Background(),
				group: &entities.AccountGroup{
					ID:    274,
					User:  &entities.User{ID: 235},
					Title: "some title",
					Accounts: entities.Accounts{
						{ID: 135}, {ID: 55},
					},
				},
			},
			delay:       10 * time.Millisecond,
			secondDelay: 100 * time.Millisecond,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalln(err)
			}
			g := tt.args.group
			result := sqlmock.NewResult(tt.want, 1)
			if tt.resultError {
				result = sqlmock.NewErrorResult(errors.New("failed"))
			}
			mock.ExpectBegin()

			mock.ExpectExec("insert into account_groups").
				WithArgs(g.User.ID, g.Title).
				WillReturnResult(result).
				WillDelayFor(tt.delay)

			if len(g.Accounts) > 0 {
				args := make([]driver.Value, len(g.Accounts)*3)
				for i, acc := range g.Accounts {
					args[i*3] = g.User.ID
					args[i*3+1] = acc.ID
					args[i*3+2] = tt.want
				}

				mock.ExpectExec("insert into `groups`").
					WithArgs(args...).
					WillReturnResult(result).
					WillDelayFor(tt.secondDelay)
			}

			if tt.wantErr {
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			ag := &accountGroups{
				wConn:   db,
				rConn:   &sql.DB{},
				timeout: tt.fields.timeout,
				manager: &logger.LoggerManagerMock{},
			}
			got, err := ag.Add(tt.args.ctx, tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("accountGroups.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("accountGroups.Add() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); !tt.wantErr && err != nil {
				t.Errorf("accounts.Add() there were unfulfilled expectations: %s", err)
			}
		})
	}
}
