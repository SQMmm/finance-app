package actions

import (
	"fmt"
	"github.com/sqmmm/finance-app/actions/data"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/api/json"
	"github.com/sqmmm/finance-app/internal/logger"
	"github.com/sqmmm/finance-app/internal/middleware"
	"github.com/sqmmm/finance-app/services/add_account"
	"net/http"
)

type addAccount struct {
	manager logger.LoggerManager
	service add_account.AccountAdder
}

func NewAddAccount(m logger.LoggerManager, s add_account.AccountAdder) http.HandlerFunc {
	return addAccount{
		manager: m,
		service: s,
	}.handle
}

func (aa addAccount) handle(w http.ResponseWriter, req *http.Request) {
	log := aa.manager.LogCtx(req.Context())
	user, ok := req.Context().Value(middleware.UserKey).(*entities.User)
	if !ok {
		data.WriteResponse(log, w, http.StatusUnauthorized, data.ErrorUnauthorized)
		return
	}

	request := &data.Account{}
	if err := json.ReadData(req, request); err != nil {
		data.WriteResponse(log, w, http.StatusBadRequest, data.ErrorBadRequest(fmt.Errorf("failed to unmarshal request: %s", err)))
		return
	}

	if err := request.Validate(); err != nil {
		data.WriteResponse(log, w, http.StatusBadRequest, data.ErrorBadRequest(err))
		return
	}

	account := request.GetEntity()
	account.User = user
	account, err := aa.service.AddAccount(req.Context(), account)
	if err != nil {
		log.Errorf("failed to add service: %s", err)
		data.WriteResponse(log, w, http.StatusInternalServerError, data.ErrorInternal)
		return
	}

	data.WriteResponse(log, w, http.StatusOK, data.Result(data.GetAccountFromEntity(account)))
}
