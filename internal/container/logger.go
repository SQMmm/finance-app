package container

import (
	"fmt"
	"os"

	"github.com/sqmmm/finance-app/internal/config"
	"github.com/sqmmm/finance-app/internal/logger"
	"github.com/sqmmm/finance-app/internal/tracker"
)

type loggerManager interface {
	logger.LoggerManager
}

var manager loggerManager

func buildTrackerManager(_ *config.Config) error {
	err := tracker.InitLog(os.Stderr)
	if err != nil {
		return fmt.Errorf("log initialization err: %s", err)
	}

	manager = tracker.NewTrackerManager()

	return nil
}
