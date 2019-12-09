package container

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sqmmm/finance-app/internal/api"
	"github.com/sqmmm/finance-app/internal/config"
	"github.com/sqmmm/finance-app/internal/middleware"
	"github.com/sqmmm/finance-app/internal/server"
	"net/http"
)

var httpServer server.Server

func GetHTTPServer() server.Server {
	return httpServer
}

func buildHTTPServer(cfg *config.Config) error {
	s := api.NewHandler(manager, cfg.Listen, notFoundAction)

	err := s.RegisterHandler(api.Handler{
		Action:  nil,
		Path:    "/example",
		Methods: []string{http.MethodGet},
		Middlewares: []mux.MiddlewareFunc{
			middleware.TrackerMiddleware(manager),
			middleware.LoggingMiddleware(manager.Log()),
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to handle healthCheck")
	}
	//
	//err = s.RegisterHandlerWithAuth(api.Handler{
	//	Action:  getFeedAction,
	//	Path:    "/questions/feed",
	//	Methods: []string{http.MethodGet},
	//	Middlewares: []mux.MiddlewareFunc{
	//		middleware.Auth(manager, authChecker),
	//		middleware.ForbidNotRegistered(),
	//	},
	//})
	//if err != nil {
	//	return errors.Wrap(err, "failed to handle getFeed")
	//}
	//
	//err = s.RegisterHandlerWithAuth(api.Handler{
	//	Action:  getQuestion,
	//	Path:    "/questions/{questionID:[1-9][0-9]*}",
	//	Methods: []string{http.MethodGet},
	//	Middlewares: []mux.MiddlewareFunc{
	//		middleware.Auth(manager, authChecker),
	//		middleware.ForbidNotRegistered(),
	//	},
	//})
	//if err != nil {
	//	return errors.Wrap(err, "failed to handle GetFeed")
	//}
	//
	//err = s.RegisterHandlerWithAuth(api.Handler{
	//	Action:  addQuestion,
	//	Path:    "/questions",
	//	Methods: []string{http.MethodPost},
	//	Middlewares: []mux.MiddlewareFunc{
	//		middleware.Auth(manager, authChecker),
	//		middleware.ForbidBanned(),
	//		middleware.ForbidNotRegistered(),
	//	},
	//})
	//if err != nil {
	//	return errors.Wrap(err, "failed to handle SetRead")
	//}
	//
	//err = s.RegisterHandlerWithAuth(api.Handler{
	//	Action:  addAnswer,
	//	Path:    "/questions/{questionID:[1-9][0-9]*}/answers",
	//	Methods: []string{http.MethodPost},
	//	Middlewares: []mux.MiddlewareFunc{
	//		middleware.Auth(manager, authChecker),
	//		middleware.ForbidNotRegistered(),
	//		middleware.ForbidBanned(),
	//	},
	//})
	//if err != nil {
	//	return errors.Wrap(err, "failed to handle SetViewed")
	//}
	//
	//err = s.RegisterHandlerWithAuth(api.Handler{
	//	Action:  sendReportAction,
	//	Path:    "/sendReport",
	//	Methods: []string{http.MethodPost},
	//	Middlewares: []mux.MiddlewareFunc{
	//		middleware.Auth(manager, authChecker),
	//		middleware.ForbidNotRegistered(),
	//	},
	//})
	//if err != nil {
	//	return errors.Wrap(err, "failed to handle sendReport")
	//}
	//
	//err = s.RegisterHandlerWithAuth(api.Handler{
	//	Action:  voteAction.ServeHTTP,
	//	Path:    "/questions/{questionID:[1-9][0-9]*}/vote",
	//	Methods: []string{http.MethodPost},
	//	Middlewares: []mux.MiddlewareFunc{
	//		middleware.Auth(manager, authChecker),
	//		middleware.ForbidNotRegistered(),
	//		middleware.ForbidBanned(),
	//	},
	//})
	//if err != nil {
	//	return errors.Wrap(err, "failed to handle vote")
	//}
	//
	//err = s.RegisterHandlerWithAuth(api.Handler{
	//	Action:  getSuggestionsAction.ServeHTTP,
	//	Path:    "/categories/suggestions",
	//	Methods: []string{http.MethodGet},
	//	Middlewares: []mux.MiddlewareFunc{
	//		middleware.Auth(manager, authChecker),
	//		middleware.ForbidNotRegistered(),
	//	},
	//})
	//if err != nil {
	//	return errors.Wrap(err, "failed to handle vote")
	//}

	httpServer = s
	return nil
}
