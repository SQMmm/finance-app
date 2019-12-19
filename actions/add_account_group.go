package actions

import (
	"fmt"
	"github.com/sqmmm/finance-app/actions/data"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/api/json"
	"github.com/sqmmm/finance-app/internal/logger"
	"github.com/sqmmm/finance-app/internal/middleware"
	"github.com/sqmmm/finance-app/services/add_account_group"
	"net/http"
)

type addAccountGroup struct {
	manager logger.LoggerManager
	service add_account_group.AccountGroupAdder
}

func NewAddAccountGroup(m logger.LoggerManager, s add_account_group.AccountGroupAdder) http.HandlerFunc {
	return addAccountGroup{
		manager: m,
		service: s,
	}.handle
}

func (aac addAccountGroup) handle(w http.ResponseWriter, req *http.Request) {
	log := aac.manager.LogCtx(req.Context())
	user, ok := req.Context().Value(middleware.UserKey).(*entities.User)
	if !ok {
		data.WriteResponse(log, w, http.StatusUnauthorized, data.ErrorUnauthorized)
		return
	}

	request := &data.AccountGroup{}
	if err := json.ReadData(req, request); err != nil {
		data.WriteResponse(log, w, http.StatusBadRequest, data.ErrorBadRequest(fmt.Errorf("failed to unmarshal request: %s", err)))
		return
	}

	if err := request.Validate(); err != nil {
		data.WriteResponse(log, w, http.StatusBadRequest, data.ErrorBadRequest(err))
		return
	}

	group := request.GetEntity()
	group.User = user

	group, err := aac.service.AddAccountGroup(req.Context(), group)
	if err != nil {
		log.Errorf("failed to add category: %s", err)
		data.WriteResponse(log, w, http.StatusInternalServerError, data.ErrorInternal)
		return
	}

	data.WriteResponse(log, w, http.StatusOK, data.Result(data.GetAccountGroupFromEntity(group)))
}
