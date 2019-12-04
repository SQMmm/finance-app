package api

import (
	"context"
	"fmt"
	"github.com/sqmmm/finance-app/internal/logger"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	//"gitlab.corp.mail.ru/otvetmailru/appotvet-main/internal/logger"
	//"gitlab.corp.mail.ru/otvetmailru/appotvet-main/internal/metrics"
	"gitlab.corp.mail.ru/otvetmailru/http-server"
)

type Handler struct {
	Action      http.HandlerFunc
	Path        string
	PathPrefix  string
	Methods     []string
	Middlewares []mux.MiddlewareFunc
}

type handler struct {
	trackerManager logger.LoggerManager
	addr           string

	handlers []Handler
}

func NewHandler(m logger.LoggerManager, addr string) *handler {
	return &handler{
		trackerManager: m,
		addr:           addr,
		handlers:       make([]Handler, 0),
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
	// TODO: Disable in production
	xhttp.EnableDebugMode()
	xhttp.RegisterErrorLogger(func(message string) {
		h.trackerManager.Log().Errorf("failed to register error logger: %s", message)
	})
	//todo: to config
	xhttp.SetVerboseLevel(xhttp.DetailedErrors)

	ln, err := net.Listen("tcp", h.addr)
	if err != nil {
		return err
	}
	go h.stop(ctx, ln)

	h.registerRoutes()

	return xhttp.Serve(ln)
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
