package container

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sqmmm/finance-app/actions"
	"github.com/sqmmm/finance-app/internal/api"
	"github.com/sqmmm/finance-app/internal/middleware"
	"github.com/sqmmm/finance-app/internal/server"
	"github.com/sqmmm/finance-app/repository"
	"github.com/sqmmm/finance-app/services/add_account"
	"github.com/sqmmm/finance-app/services/add_category"
	"net/http"
	"time"

	"github.com/sqmmm/finance-app/internal/config"
)

var httpServer server.Server

// Build start dependency build process
func Build(cfg *config.Config) error {
	//init infrastructure
	tracker, err := buildTrackerManager(cfg)
	if err != nil {
		return fmt.Errorf("failed to build tracker: %s", err)
	}
	w, r, err := buildMySQLClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create mysql connections: %s", err)
	}

	//init repositories
	//todo: CONFIG
	accountRepository := repository.NewAccounts(w, r, time.Second)
	categoryRepository := repository.NewCategories(w, r, time.Second)
	iconRepository := repository.NewIcons(w, r, time.Second)

	//init services
	addAccount := add_account.NewService(tracker, accountRepository)
	addCategory := add_category.NewService(tracker, categoryRepository, iconRepository)

	//init handlers
	notFoundHandler := actions.NewNotFound(tracker)
	addAccountHandler := actions.NewAddAccount(tracker, addAccount)
	addCategoryHandler := actions.NewAddCategory(tracker, addCategory)

	//init server
	s := api.NewHandler(tracker, cfg.Listen, notFoundHandler)

	err = s.RegisterHandler(api.Handler{
		Action:  nil,
		Path:    "/example",
		Methods: []string{http.MethodGet},
		Middlewares: []mux.MiddlewareFunc{
			middleware.Tracker(tracker),
			middleware.Logging(tracker.Log()),
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to handle healthCheck")
	}
	err = s.RegisterHandler(api.Handler{
		Action:  addAccountHandler,
		Path:    "/accounts",
		Methods: []string{http.MethodPost},
		Middlewares: []mux.MiddlewareFunc{
			middleware.Tracker(tracker),
			middleware.Logging(tracker.Log()),
			middleware.Auth(tracker),
		},
		PathPrefix: "/api/v1",
	})
	if err != nil {
		return errors.Wrap(err, "failed to handle addAccount")
	}
	err = s.RegisterHandler(api.Handler{
		Action:  addCategoryHandler,
		Path:    "/categories",
		Methods: []string{http.MethodPost},
		Middlewares: []mux.MiddlewareFunc{
			middleware.Tracker(tracker),
			middleware.Logging(tracker.Log()),
			middleware.Auth(tracker),
		},
		PathPrefix: "/api/v1",
	})
	if err != nil {
		return errors.Wrap(err, "failed to handle addCategory")
	}

	httpServer = s

	return nil
}

func GetServer() server.Server {
	return httpServer
}
