package actions

import (
	"bytes"
	"context"
	"errors"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sqmmm/finance-app/internal/logger"
)

func Test_getAccount_handle(t *testing.T) {
	type data struct {
		addAccountData
	}
	type mocks struct {
		addAccountMock
	}
	tests := []struct {
		name string
		data
		mocks
		response string
		request  string
		wantCode int
		user     *entities.User
	}{
		{
			name:    "success",
			request: `{"title":"some title","start_balance":89.25,"use_in_reports":true}`,
			user:    &entities.User{ID: 143},
			data: data{
				addAccountData: addAccountData{account: entities.Account{
					User:         &entities.User{ID: 143},
					Title:        "some title",
					StartBalance: 89.25,
					UseInReports: true,
				}},
			},
			mocks: mocks{addAccountMock: addAccountMock{
				acc: &entities.Account{
					ID:           1234,
					User:         &entities.User{ID: 143},
					Title:        "some title",
					StartBalance: 89.25,
					UseInReports: true,
				},
			}},
			response: `{"result":{"id":1234,"title":"some title","user_id":143,"start_balance":89.25,"use_in_reports":true}}
`,
			wantCode: http.StatusOK,
		},
		{
			name:    "auth error",
			request: `{"title":"some title","start_balance":89.25,"use_in_reports":true}`,
			response: `{"error":{"code":401,"message":"authorization required for this action"}}
`,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:    "unmarshal error",
			request: `efdef`,
			user:    &entities.User{ID: 143},
			response: `{"error":{"code":400,"message":"failed to unmarshal request: invalid character 'e' looking for beginning of value"}}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "validation error",
			request: `{}`,
			user:    &entities.User{ID: 143},
			response: `{"error":{"code":400,"message":"title length is incorrect: cannot be blank."}}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "service error",
			request: `{"title":"some title","start_balance":89.25,"use_in_reports":true}`,
			user:    &entities.User{ID: 143},
			data: data{
				addAccountData: addAccountData{account: entities.Account{
					User:         &entities.User{ID: 143},
					Title:        "some title",
					StartBalance: 89.25,
					UseInReports: true,
				}},
			},
			mocks: mocks{addAccountMock: addAccountMock{
				err: errors.New("failed"),
			}},
			response: `{"error":{"code":500,"message":"internal error"}}
`,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := "/api/v1/accounts"
			buf := bytes.NewBufferString(tt.request)
			req := httptest.NewRequest(http.MethodPost, target, buf)
			if tt.user != nil {
				req = req.WithContext(context.WithValue(context.Background(), middleware.UserKey, tt.user))
			}
			rr := httptest.NewRecorder()

			service := &accountAdderMock{
				addAccountMock: tt.mocks.addAccountMock,
			}
			fn := NewAddAccount(&logger.LoggerManagerMock{}, service)
			fn(rr, req)

			if code := rr.Code; code != tt.wantCode {
				t.Errorf("handler.AddAccount() must response with status=%v, but responsed=%v", tt.wantCode, code)
			}
			assert.Equal(t, tt.response, rr.Body.String(), "checking response data")
			assert.Equal(t, tt.data.addAccountData, service.addAccountData, "checking add account data")
		})
	}
}

type addAccountData struct {
	account entities.Account
}
type addAccountMock struct {
	acc *entities.Account
	err error
}
type accountAdderMock struct {
	addAccountData

	addAccountMock
}

func (aam *accountAdderMock) AddAccount(_ context.Context, acc *entities.Account) (*entities.Account, error) {
	aam.addAccountData = addAccountData{account: *acc}
	return aam.addAccountMock.acc, aam.addAccountMock.err
}
