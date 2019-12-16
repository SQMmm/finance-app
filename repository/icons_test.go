package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/sqmmm/finance-app/entities"
)

func Test_icons_GetIconByID(t *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     *entities.Icon
		wantErr  bool
		delay    time.Duration
		rowsMock [][]driver.Value
	}{
		{
			name:   "success",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				id:  163,
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{124, "name", "path"},
				{125, "name2", "path2"},
			},
			want: &entities.Icon{ID: 124, Name: "name", Path: "path"},
		},
		{
			name:   "failed by timeout",
			fields: fields{timeout: time.Millisecond},
			args: args{
				ctx: context.Background(),
				id:  163,
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{124, "name", "path"},
				{125, "name2", "path2"},
			},
			wantErr: true,
		},
		{
			name:   "failed by canceled context",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: ctxCanceled,
				id:  163,
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{124, "name", "path"},
				{125, "name2", "path2"},
			},
			wantErr: true,
		},
		{
			name:   "failed by scan error",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				id:  163,
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{nil, "name", "path"},
				{125, "name2", "path2"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		db, mock, err := sqlmock.New()
		if err != nil {
			log.Fatalln(err)
		}

		rows := sqlmock.NewRows([]string{"id", "name", "path"})
		for _, m := range tt.rowsMock {
			rows = rows.AddRow(m...)
		}

		mock.ExpectQuery("select id, name, path from icons where id = ?").
			WithArgs(tt.args.id).
			RowsWillBeClosed().
			WillDelayFor(tt.delay).
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			i := &icons{
				wConn:   &sql.DB{},
				rConn:   db,
				timeout: tt.fields.timeout,
			}
			got, err := i.GetIconByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("icons.GetIconByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("icons.GetIconByID() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); !tt.wantErr && err != nil {
				t.Errorf("icons.GetIconByID() there were unfulfilled expectations: %s", err)
			}
		})
	}
}
func Test_icons_GetColorByID(t *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     *entities.IconColor
		wantErr  bool
		delay    time.Duration
		rowsMock [][]driver.Value
	}{
		{
			name:   "success",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				id:  163,
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{124, "name", "color"},
				{125, "name2", "color2"},
			},
			want: &entities.IconColor{ID: 124, Name: "name", Color: "color"},
		},
		{
			name:   "failed by timeout",
			fields: fields{timeout: time.Millisecond},
			args: args{
				ctx: context.Background(),
				id:  163,
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{124, "name", "color"},
			},
			wantErr: true,
		},
		{
			name:   "failed by canceled context",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: ctxCanceled,
				id:  163,
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{124, "name", "color"},
			},
			wantErr: true,
		},
		{
			name:   "failed by scan error",
			fields: fields{timeout: time.Second},
			args: args{
				ctx: context.Background(),
				id:  163,
			},
			delay: 10 * time.Millisecond,
			rowsMock: [][]driver.Value{
				{nil, "name", "color"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		db, mock, err := sqlmock.New()
		if err != nil {
			log.Fatalln(err)
		}

		rows := sqlmock.NewRows([]string{"id", "name", "color"})
		for _, m := range tt.rowsMock {
			rows = rows.AddRow(m...)
		}

		mock.ExpectQuery("select id, name, color from icon_colors where id = ?").
			WithArgs(tt.args.id).
			RowsWillBeClosed().
			WillDelayFor(tt.delay).
			WillReturnRows(rows)

		t.Run(tt.name, func(t *testing.T) {
			i := &icons{
				wConn:   &sql.DB{},
				rConn:   db,
				timeout: tt.fields.timeout,
			}
			got, err := i.GetColorByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("icons.GetColorByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("icons.GetColorByID() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); !tt.wantErr && err != nil {
				t.Errorf("icons.GetColorByID() there were unfulfilled expectations: %s", err)
			}
		})
	}
}
