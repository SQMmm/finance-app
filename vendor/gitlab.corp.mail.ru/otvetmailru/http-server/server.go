package xhttp

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	router        = mux.NewRouter()
	defaultServer = &http.Server{
		Handler: router,
	}
)

func Router() *mux.Router {
	return router
}

// Serve start serving application defaultServer
func Serve(ln net.Listener) error {
	return defaultServer.Serve(ln)
}
