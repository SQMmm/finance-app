package api

import (
	"net/http"

	//"gitlab.corp.mail.ru/otvetmailru/appotvet-main/internal/middleware"
	"gitlab.corp.mail.ru/otvetmailru/http-server"
)

func (h *handler) registerRoutes() {
	router := xhttp.Router()

	for _, handler := range h.handlers {
		privateAPI := router.PathPrefix(handler.PathPrefix).Subrouter()
		for _, m := range handler.Middlewares {
			privateAPI.Use(m)
		}
		privateAPI.Methods(handler.Methods...).Path(handler.Path).HandlerFunc(handler.Action)
	}
	router.NotFoundHandler = http.HandlerFunc(h.notFoundHandler)
}

func (h *handler) notFoundHandler(w http.ResponseWriter, req *http.Request) {
	err := xhttp.WriteJsonNotFound(w)
	if err != nil {
		h.trackerManager.Log().Errorf("failed to write response")
	}
}
