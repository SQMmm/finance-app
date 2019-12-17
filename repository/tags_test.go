package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
	"time"

	"github.com/sqmmm/finance-app/entities"
)

func Test_tags_Add(t *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	type args struct {
		ctx context.Context
		tag *entities.Tag
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        int64
		wantErr     bool
		delay       time.Duration
		resultError bool
	}{
		{
			name:   "success",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				tag: &entities.Tag{
					ID:    635,
					Title: "tag title",
					User:  &entities.User{ID: 11},
				},
			},
			delay: 10 * time.Millisecond,
			want:  168,
		},
		{
			name:   "failed by timeout",
			fields: fields{timeout: time.Millisecond},
			args: args{
				ctx: context.Background(),
				tag: &entities.Tag{
					ID:    635,
					Title: "tag title",
					User:  &entities.User{ID: 11},
				},
			},
			delay:   10 * time.Millisecond,
			wantErr: true,
		},
		{
			name:   "failed by canceled context",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: ctxCanceled,
				tag: &entities.Tag{
					ID:    635,
					Title: "tag title",
					User:  &entities.User{ID: 11},
				},
			},
			delay:   10 * time.Millisecond,
			wantErr: true,
		},
		{
			name:   "failed by by getting last insert id",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				tag: &entities.Tag{
					ID:    635,
					Title: "tag title",
					User:  &entities.User{ID: 11},
				},
			},
			delay:       10 * time.Millisecond,
			resultError: true,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalln(err)
			}
			result := sqlmock.NewResult(tt.want, 1)
			if tt.resultError {
				result = sqlmock.NewErrorResult(errors.New("failed"))
			}

			mock.ExpectExec("insert into tags").
				WithArgs(tt.args.tag.User.ID, tt.args.tag.Title).
				WillReturnResult(result).
				WillDelayFor(tt.delay)

			tg := &tags{
				wConn:   db,
				rConn:   &sql.DB{},
				timeout: tt.fields.timeout,
			}
			got, err := tg.Add(tt.args.ctx, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("tags.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tags.Add() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); !tt.wantErr && err != nil {
				t.Errorf("tags.Add() there were unfulfilled expectations: %s", err)
			}
		})
	}
}
