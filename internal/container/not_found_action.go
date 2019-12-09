package container

import (
	"github.com/sqmmm/finance-app/actions"
	"github.com/sqmmm/finance-app/internal/config"
	"net/http"
)

var notFoundAction http.HandlerFunc

func buildNotFoundAction(cfg *config.Config) error {
	notFoundAction = actions.NewNotFound(manager)

	return nil
}
