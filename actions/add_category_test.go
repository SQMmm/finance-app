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

func Test_addCategory_handle(t *testing.T) {
	type data struct {
		addCategoryData
	}
	type mocks struct {
		addCategoryMock
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
			name:    "success with icon and color",
			request: `{"title":"some title","icon":{"id":123},"color":{"id":11}}`,
			user:    &entities.User{ID: 7267},
			data: data{addCategoryData: addCategoryData{category: entities.Category{
				Title: "some title",
				User:  &entities.User{ID: 7267},
				Icon:  &entities.Icon{ID: 123},
				Color: &entities.IconColor{ID: 11},
			}}},
			mocks: mocks{addCategoryMock: addCategoryMock{category: &entities.Category{
				ID:    75,
				Title: "some title",
				User:  &entities.User{ID: 7267},
				Icon:  &entities.Icon{ID: 123, Path: "icon path", Name: "icon name"},
				Color: &entities.IconColor{ID: 11, Color: "icon color", Name: "color name"},
			}}},
			response: `{"result":{"id":75,"user_id":7267,"title":"some title","icon":{"id":123,"name":"icon name",` +
				`"path":"icon path"},"color":{"id":11,"name":"color name","color":"icon color"}}}
`,
			wantCode: http.StatusOK,
		},
		{
			name:    "success with icon and color",
			request: `{"title":"some title"}`,
			user:    &entities.User{ID: 7267},
			data: data{addCategoryData: addCategoryData{category: entities.Category{
				Title: "some title",
				User:  &entities.User{ID: 7267},
			}}},
			mocks: mocks{addCategoryMock: addCategoryMock{category: &entities.Category{
				ID:    75,
				Title: "some title",
				User:  &entities.User{ID: 7267},
			}}},
			response: `{"result":{"id":75,"user_id":7267,"title":"some title"}}
`,
			wantCode: http.StatusOK,
		},
		{
			name:    "failed to get user",
			request: `{"title":"some title","icon_id":123,"color_id":11}`,
			response: `{"error":{"code":401,"message":"authorization required for this action"}}
`,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:    "failed to unmarshal data",
			request: `{fefefefef}`,
			user:    &entities.User{ID: 7267},
			response: `{"error":{"code":400,"message":"failed to unmarshal request: invalid character 'f' looking for beginning of object key string"}}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "failed to validate data",
			request: `{}`,
			user:    &entities.User{ID: 7267},
			response: `{"error":{"code":400,"message":"title length is incorrect: cannot be blank."}}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "failed to add category",
			request: `{"title":"some title","icon":{"id":123},"color":{"id":11}}`,
			user:    &entities.User{ID: 7267},
			data: data{addCategoryData: addCategoryData{category: entities.Category{
				Title: "some title",
				User:  &entities.User{ID: 7267},
				Icon:  &entities.Icon{ID: 123},
				Color: &entities.IconColor{ID: 11},
			}}},
			mocks: mocks{addCategoryMock: addCategoryMock{err: errors.New("failed")}},
			response: `{"error":{"code":500,"message":"internal error"}}
`,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := "/api/v1/categories"
			buf := bytes.NewBufferString(tt.request)
			req := httptest.NewRequest(http.MethodPost, target, buf)
			if tt.user != nil {
				req = req.WithContext(context.WithValue(context.Background(), middleware.UserKey, tt.user))
			}
			rr := httptest.NewRecorder()

			service := &categoryAdderMock{
				addCategoryMock: tt.mocks.addCategoryMock,
			}
			fn := NewAddCategory(&logger.LoggerManagerMock{}, service)
			fn(rr, req)

			if code := rr.Code; code != tt.wantCode {
				t.Errorf("handler.AddAccount() must response with status=%v, but responsed=%v", tt.wantCode, code)
			}
			assert.Equal(t, tt.response, rr.Body.String(), "checking response data")
			assert.Equal(t, tt.data.addCategoryData, service.addCategoryData, "checking add category data")
		})
	}
}

type addCategoryData struct {
	category entities.Category
}
type addCategoryMock struct {
	category *entities.Category
	err      error
}
type categoryAdderMock struct {
	addCategoryData

	addCategoryMock
}

func (cam *categoryAdderMock) AddCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	cam.addCategoryData = addCategoryData{category: *category}
	return cam.addCategoryMock.category, cam.addCategoryMock.err
}
