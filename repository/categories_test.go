package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
	"time"

	"github.com/sqmmm/finance-app/entities"
)

func Test_categories_Add(t *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	type args struct {
		ctx      context.Context
		category *entities.Category
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
			name:   "success without icon and color",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				category: &entities.Category{
					ID:    124,
					Title: "some title",
					User:  &entities.User{ID: 52},
				},
			},
			delay: 10 * time.Millisecond,
			want:  736,
		},
		{
			name:   "success with icon and color",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				category: &entities.Category{
					ID:    124,
					Title: "some title",
					User:  &entities.User{ID: 52},
					Icon:  &entities.Icon{ID: 535},
					Color: &entities.IconColor{ID: 11},
				},
			},
			delay: 10 * time.Millisecond,
			want:  736,
		},
		{
			name:   "failed by timeout",
			fields: fields{timeout: time.Millisecond},
			args: args{
				ctx: context.Background(),
				category: &entities.Category{
					ID:    124,
					Title: "some title",
					User:  &entities.User{ID: 52},
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
				category: &entities.Category{
					ID:    124,
					Title: "some title",
					User:  &entities.User{ID: 52},
				},
			},
			delay:   10 * time.Millisecond,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalln(err)
			}
			cat := tt.args.category
			args := []driver.Value{cat.User.ID, cat.Title}
			if cat.Icon != nil {
				args = append(args, cat.Icon.ID)
			}
			if cat.Color != nil {
				args = append(args, cat.Color.ID)
			}

			mock.ExpectExec("insert into categories").
				WithArgs(args...).
				WillReturnResult(sqlmock.NewResult(tt.want, 1)).
				WillDelayFor(tt.delay)

			c := NewCategories(db, &sql.DB{}, tt.fields.timeout)
			got, err := c.Add(tt.args.ctx, tt.args.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("categories.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("categories.Add() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); !tt.wantErr && err != nil {
				t.Errorf("categories.Add() there were unfulfilled expectations: %s", err)
			}
		})
	}
}
