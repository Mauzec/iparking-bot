package main

import (
	"log"

	"github.com/mauzec/ibot-things/config"
	"github.com/mauzec/ibot-things/internal/bot"
	"github.com/mauzec/ibot-things/pkg/logger"
)

// TODO: use logger.logger in all functions in all project
func main() {
	if err := logger.Init(true); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	cfg, err := config.LoadConfig("app", "env", "config")
	if err != nil {
		logger.Logger().Fatalf("failed to load config: %v", err)
	}

	b, err := bot.NewBotFromConfig(&cfg)
	if err != nil {
		logger.Logger().Fatalf("failed to create bot: %v", err)
	}

	b.Run()
}
