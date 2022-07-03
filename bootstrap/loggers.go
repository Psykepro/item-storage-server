package bootstrap

import (
	domain "github.com/Psykepro/item-storage-server/_domain"
	"github.com/Psykepro/item-storage-server/config"
	"github.com/Psykepro/item-storage-server/pkg/logger"
)

func Loggers(cfg *config.Config) (domain.Logger, domain.Logger) {
	stdOutLogger := logger.NewLogger(cfg.Logging)
	stdOutLogger.InitStdOutLogger()
	stdOutLogger.Infof("LogLevel: %s, Mode: %s", cfg.Logging.StdOut.Level, cfg.Logging.StdOut.Mode)
	fileLogger := logger.NewLogger(cfg.Logging)
	fileLogger.InitFileLogger()

	return stdOutLogger, fileLogger
}
