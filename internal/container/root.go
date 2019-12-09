package container

import (
	"github.com/pkg/errors"

	"github.com/sqmmm/finance-app/internal/config"
)

// Build start dependency build process
func Build(cfg *config.Config) error {
	builders := []func(*config.Config) error{
		//infra
		buildTrackerManager,

		//actions
		buildNotFoundAction,

		//servers
		buildHTTPServer,
	}

	for _, builder := range builders {
		if err := builder(cfg); err != nil {
			return errors.Wrap(err, "unable to start application")
		}
	}

	return nil
}
