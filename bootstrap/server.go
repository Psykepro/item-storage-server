package bootstrap

import (
	domain "github.com/Psykepro/item-storage-server/_domain"
	"github.com/Psykepro/item-storage-server/config"
	"github.com/Psykepro/item-storage-server/internal/item"
	"github.com/Psykepro/item-storage-server/internal/server"
)

func Server(stdOutLogger domain.Logger, fileLogger domain.Logger, cfg *config.Config) *server.Server {
	itemService := item.NewService(stdOutLogger, fileLogger)
	itemRequestHandler := item.NewRequestHandler(itemService, stdOutLogger)
	s := server.NewServer(cfg.RabbitMQ, itemRequestHandler, stdOutLogger)
	return s
}
