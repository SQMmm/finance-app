package actions

import (
	"fmt"
	"github.com/sqmmm/finance-app/actions/data"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/api/json"
	"github.com/sqmmm/finance-app/internal/logger"
	"github.com/sqmmm/finance-app/internal/middleware"
	"github.com/sqmmm/finance-app/services/add_category"
	"net/http"
)

type addCategory struct {
	manager logger.LoggerManager
	service add_category.CategoryAdder
}

func NewAddCategory(m logger.LoggerManager, s add_category.CategoryAdder) http.HandlerFunc {
	return addCategory{
		manager: m,
		service: s,
	}.handle
}

func (ac addCategory) handle(w http.ResponseWriter, req *http.Request) {
	log := ac.manager.LogCtx(req.Context())
	user, ok := req.Context().Value(middleware.UserKey).(*entities.User)
	if !ok {
		data.WriteResponse(log, w, http.StatusUnauthorized, data.ErrorUnauthorized)
		return
	}

	request := &data.Category{}
	if err := json.ReadData(req, request); err != nil {
		data.WriteResponse(log, w, http.StatusBadRequest, data.ErrorBadRequest(fmt.Errorf("failed to unmarshal request: %s", err)))
		return
	}

	if err := request.Validate(); err != nil {
		data.WriteResponse(log, w, http.StatusBadRequest, data.ErrorBadRequest(err))
		return
	}

	cat := request.GetEntity()
	cat.User = user

	cat, err := ac.service.AddCategory(req.Context(), cat)
	if err != nil {
		log.Errorf("failed to add category: %s", err)
		data.WriteResponse(log, w, http.StatusInternalServerError, data.ErrorInternal)
		return
	}

	data.WriteResponse(log, w, http.StatusOK, data.Result(data.GetCategoryFromEntity(cat)))
}
