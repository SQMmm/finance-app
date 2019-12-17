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

func Test_addTag_handle(t *testing.T) {
	type data struct {
		addTagData
	}
	type mocks struct {
		addTagMock
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
			request: `{"title":"tag title"}`,
			data: data{addTagData: addTagData{tag: entities.Tag{
				Title: "tag title",
				User: &entities.User{ID: 725},
			}}},
			mocks: mocks{addTagMock: addTagMock{
				tag:&entities.Tag{ID:256, User:&entities.User{ID:736}, Title:"some title"},
			}},
			user:&entities.User{ID:725},
			response: `{"result":{"id":256,"user_id":736,"title":"some title"}}
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
			request: `{"title":"tag title"}`,
			data: data{addTagData: addTagData{tag: entities.Tag{
				Title: "tag title",
				User: &entities.User{ID: 725},
			}}},
			mocks: mocks{addTagMock: addTagMock{
				err: errors.New("failed"),
			}},
			user:&entities.User{ID:725},
			response: `{"error":{"code":500,"message":"internal error"}}
`,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := "/api/v1/tags"
			buf := bytes.NewBufferString(tt.request)
			req := httptest.NewRequest(http.MethodPost, target, buf)
			if tt.user != nil {
				req = req.WithContext(context.WithValue(context.Background(), middleware.UserKey, tt.user))
			}
			rr := httptest.NewRecorder()
			s := &tagAdderMock{
				addTagMock: tt.mocks.addTagMock,
			}
			fn := NewAddTagHandler(&logger.LoggerManagerMock{}, s)
			fn(rr, req)
			if code := rr.Code; code != tt.wantCode {
				t.Errorf("handler.AddTag() must response with status=%v, but responsed=%v", tt.wantCode, code)
			}
			assert.Equal(t, tt.response, rr.Body.String(), "checking response data")
			assert.Equal(t, tt.data.addTagData, s.addTagData, "checking add tag data")
		})
	}
}

type addTagData struct {
	tag entities.Tag
}
type addTagMock struct {
	tag *entities.Tag
	err error
}
type tagAdderMock struct {
	addTagData

	addTagMock
}

func (tam *tagAdderMock) AddTag(_ context.Context, tag *entities.Tag) (*entities.Tag, error) {
	tam.addTagData = addTagData{tag: *tag}
	return tam.addTagMock.tag, tam.addTagMock.err
}
