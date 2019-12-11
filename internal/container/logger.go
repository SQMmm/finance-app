package container

import (
	"fmt"
	"os"

	"github.com/sqmmm/finance-app/internal/config"
	"github.com/sqmmm/finance-app/internal/logger"
	"github.com/sqmmm/finance-app/internal/tracker"
)

func buildTrackerManager(_ *config.Config) (logger.LoggerManager, error) {
	err := tracker.InitLog(os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("log initialization err: %s", err)
	}

	return tracker.NewTrackerManager(),nil
}
