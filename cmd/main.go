package main

import (
	"github.com/Psykepro/item-storage-server/bootstrap"
)

func main() {
	cfg := bootstrap.Config()
	stdOutLogger, fileLogger := bootstrap.Loggers(cfg)
	s := bootstrap.Server(stdOutLogger, fileLogger, cfg)

	s.Run()
}
