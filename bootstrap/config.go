package bootstrap

import (
	"log"

	"github.com/Psykepro/item-storage-server/config"
)

func Config() *config.Config {
	log.Println("Starting server")
	cfg, err := config.GetConfig(config.Path)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}
	log.Println("Successfully loaded config.")

	return cfg
}
