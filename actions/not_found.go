package actions

import (
	"github.com/sqmmm/finance-app/actions/data"
	"github.com/sqmmm/finance-app/internal/api/json"
	"github.com/sqmmm/finance-app/internal/logger"
	"net/http"
)

type notFound struct {
	tracker logger.LoggerManager
}

func NewNotFound(tr logger.LoggerManager) http.HandlerFunc {
	return notFound{
		tracker: tr,
	}.notFoundHandler
}

func (nf notFound) notFoundHandler(w http.ResponseWriter, req *http.Request) {
	err := json.WriteResponse(w, http.StatusNotFound, data.ErrorNotFound)
	if err != nil {
		nf.tracker.LogCtx(req.Context()).Errorf("failed to write response: %s", err)
	}
}
