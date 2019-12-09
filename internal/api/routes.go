package api

import (
	"net/http"
)

func (h *handler) registerRoutes() {
	router := h.router

	for _, handler := range h.handlers {
		privateAPI := router.PathPrefix(handler.PathPrefix).Subrouter()
		for _, m := range handler.Middlewares {
			privateAPI.Use(m)
		}
		privateAPI.Methods(handler.Methods...).Path(handler.Path).HandlerFunc(handler.Action)
	}
	router.NotFoundHandler = http.HandlerFunc(h.notFoundHandler)
}
