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

func Test_addAccountGroup_handle(t *testing.T) {
	type data struct {
		addAccountGroupData
	}
	type mocks struct {
		addAccountGroupMock
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
			name:     "success with accounts",
			wantCode: http.StatusOK,
			request:  `{"title":"group title","accounts":[{"id":1},{"id":2}]}`,
			data: data{addAccountGroupData: addAccountGroupData{group: entities.AccountGroup{
				Title:    "group title",
				Accounts: entities.Accounts{{ID: 1}, {ID: 2}},
				User:     &entities.User{ID: 554},
			}}},
			mocks: mocks{addAccountGroupMock: addAccountGroupMock{
				group: &entities.AccountGroup{
					ID:    134,
					Title: "group title",
					Accounts: entities.Accounts{
						{ID: 1, Title: "1dwd", User: &entities.User{ID: 111}, UseInReports: true},
						{ID: 2, StartBalance: 245, User: &entities.User{ID: 111}},
					},
					User: &entities.User{ID: 554},
				},
			}},
			user: &entities.User{ID: 554},
			response: `{"result":{"id":134,"user_id":554,"title":"group title","accounts":` +
				`[{"id":1,"title":"1dwd","user_id":111,"start_balance":0,"use_in_reports":true},` +
				`{"id":2,"title":"","user_id":111,"start_balance":245,"use_in_reports":false}]}}
`,
		},
		{
			name:     "success without accounts",
			wantCode: http.StatusOK,
			request:  `{"title":"group title"}`,
			data: data{addAccountGroupData: addAccountGroupData{group: entities.AccountGroup{
				Title:    "group title",
				User:     &entities.User{ID: 554},
				Accounts: entities.Accounts{},
			}}},
			mocks: mocks{addAccountGroupMock: addAccountGroupMock{
				group: &entities.AccountGroup{
					ID:    134,
					Title: "group title",
					User:  &entities.User{ID: 554},
				},
			}},
			user: &entities.User{ID: 554},
			response: `{"result":{"id":134,"user_id":554,"title":"group title"}}
`,
		},
		{
			name:    "failed to get user",
			request: `{"title":"group title"}`,
			response: `{"error":{"code":401,"message":"authorization required for this action"}}
`,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "failed to unmarshal data",
			wantCode: http.StatusBadRequest,
			request:  `{ascsc}`,
			user:     &entities.User{ID: 554},
			response: `{"error":{"code":400,"message":"failed to unmarshal request: invalid character 'a' looking for beginning of object key string"}}
`,
		},
		{
			name:     "failed to validate data",
			wantCode: http.StatusBadRequest,
			request:  `{}`,
			user:     &entities.User{ID: 554},
			response: `{"error":{"code":400,"message":"title length is incorrect: cannot be blank."}}
`,
		},
		{
			name:    "failed to add grou0p",
			request: `{"title":"group title"}`,
			data: data{addAccountGroupData: addAccountGroupData{group: entities.AccountGroup{
				Title:    "group title",
				User:     &entities.User{ID: 554},
				Accounts: entities.Accounts{},
			}}},
			mocks: mocks{addAccountGroupMock: addAccountGroupMock{
				err: errors.New("failed"),
			}},
			user: &entities.User{ID: 554},
			response: `{"error":{"code":500,"message":"internal error"}}
`,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := "/api/v1/accounts/group"
			buf := bytes.NewBufferString(tt.request)
			req := httptest.NewRequest(http.MethodPost, target, buf)
			if tt.user != nil {
				req = req.WithContext(context.WithValue(context.Background(), middleware.UserKey, tt.user))
			}
			rr := httptest.NewRecorder()
			service := &accountGroupAdderMock{
				addAccountGroupMock: tt.mocks.addAccountGroupMock,
			}
			fn := NewAddAccountGroup(&logger.LoggerManagerMock{}, service)
			fn(rr, req)

			if code := rr.Code; code != tt.wantCode {
				t.Errorf("handler.AddAccountGroup() must response with status=%v, but responsed=%v", tt.wantCode, code)
			}
			assert.Equal(t, tt.response, rr.Body.String(), "checking response data")
			assert.Equal(t, tt.data.addAccountGroupData, service.addAccountGroupData, "checking add account group data")
		})
	}
}

type addAccountGroupData struct {
	group entities.AccountGroup
}
type addAccountGroupMock struct {
	group *entities.AccountGroup
	err   error
}
type accountGroupAdderMock struct {
	addAccountGroupData

	addAccountGroupMock
}

func (agam *accountGroupAdderMock) AddAccountGroup(ctx context.Context, group *entities.AccountGroup) (*entities.AccountGroup, error) {
	agam.addAccountGroupData = addAccountGroupData{group: *group}
	return agam.addAccountGroupMock.group, agam.addAccountGroupMock.err
}
