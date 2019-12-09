package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/sqmmm/finance-app/internal/logger"
)

type Handler struct {
	Action      http.HandlerFunc
	Path        string
	PathPrefix  string
	Methods     []string
	Middlewares []mux.MiddlewareFunc
}

type handler struct {
	trackerManager  logger.LoggerManager
	addr            string
	notFoundHandler http.HandlerFunc

	handlers []Handler
	router   *mux.Router
	server   *http.Server
}

func NewHandler(m logger.LoggerManager, addr string, notFoundHandler http.HandlerFunc) *handler {
	r := mux.NewRouter()
	return &handler{
		trackerManager:  m,
		addr:            addr,
		notFoundHandler: notFoundHandler,
		handlers:        make([]Handler, 0),
		router:          r,
		server: &http.Server{
			Handler: r,
		},
	}
}

func (h *handler) RegisterHandler(handler Handler) error {
	for _, method := range handler.Methods {
		if err := checkMethod(method); err != nil {
			return err
		}
	}
	h.handlers = append(h.handlers, handler)

	return nil
}

// Serve start serving application api
func (h *handler) Serve(ctx context.Context, _ *sync.WaitGroup) error {
	ln, err := net.Listen("tcp", h.addr)
	if err != nil {
		return err
	}
	go h.stop(ctx, ln)

	h.registerRoutes()

	return h.server.Serve(ln)
}

func (h *handler) stop(ctx context.Context, ln net.Listener) {
	select {
	case <-ctx.Done():
		//todo: add to xhttp method server.Shutdown()
	}
}

func checkMethod(method string) error {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete,
		http.MethodConnect, http.MethodOptions, http.MethodTrace:
	default:
		return fmt.Errorf("failed to register handler: method `%s` is not supported", method)
	}

	return nil
}
