package actions

import (
	"fmt"
	"github.com/sqmmm/finance-app/actions/data"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/api/json"
	"github.com/sqmmm/finance-app/internal/logger"
	"github.com/sqmmm/finance-app/internal/middleware"
	"github.com/sqmmm/finance-app/services/add_tag"
	"net/http"
)

type addTag struct {
	manager logger.LoggerManager
	service add_tag.TagAdder
}

func NewAddTagHandler(m logger.LoggerManager, s add_tag.TagAdder) http.HandlerFunc {
	return addTag{
		manager: m,
		service: s,
	}.handle
}

func (ad addTag) handle(w http.ResponseWriter, req *http.Request) {
	log := ad.manager.LogCtx(req.Context())
	user, ok := req.Context().Value(middleware.UserKey).(*entities.User)
	if !ok {
		data.WriteResponse(log, w, http.StatusUnauthorized, data.ErrorUnauthorized)
		return
	}
	request := &data.Tag{}

	if err := json.ReadData(req, request); err != nil {
		data.WriteResponse(log, w, http.StatusBadRequest, data.ErrorBadRequest(fmt.Errorf("failed to unmarshal request: %s", err)))
		return
	}

	if err := request.Validate(); err != nil {
		data.WriteResponse(log, w, http.StatusBadRequest, data.ErrorBadRequest(err))
		return
	}

	tag := request.GetEntity()
	tag.User = user

	tag, err := ad.service.AddTag(req.Context(), tag)
	if err != nil {
		log.Errorf("failed to add tag: %s", err)
		data.WriteResponse(log, w, http.StatusInternalServerError, data.ErrorInternal)
		return
	}

	data.WriteResponse(log, w, http.StatusOK, data.Result(data.GetTagFromEntity(tag)))
}
