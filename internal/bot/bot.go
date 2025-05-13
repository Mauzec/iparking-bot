package bot

import (
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mauzec/ibot-things/config"
)

type Bot struct {
	api           *tgbotapi.BotAPI
	handler       *Handler
	updateTimeout int
}

// func NewBot(token, dataPath string) (*Bot, error) {
// 	api, err := tgbotapi.NewBotAPI(token)
// 	if err != nil {
// 		return nil, err
// 	}
// 	api.Debug = true

// 	return &Bot{
// 		api:      api,
// 		dataPath: dataPath,
// 	}, nil
// }

func NewBotFromConfig(cfg *config.Config) (*Bot, error) {
	if cfg.BotToken == "" {
		return nil, errors.New("bot token is empty")
	}
	if cfg.DataPath == "" {
		return nil, errors.New("data path is empty")
	}
	if cfg.BotTimeout <= 0 {
		cfg.BotTimeout = 60
		log.Printf("[BOT] bot timeout is not set or <= 0, using default value: %d", cfg.BotTimeout)
	}

	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	api.Debug = true

	h := NewHandler(cfg.DataPath)

	return &Bot{
		api:           api,
		handler:       h,
		updateTimeout: cfg.BotTimeout,
	}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.updateTimeout

	updates := b.api.GetUpdatesChan(u)
	for upd := range updates {
		if upd.Message != nil && upd.Message.IsCommand() {
			if err := b.handler.Handle(b.api, upd); err != nil {
				log.Printf("[BOT] handler error: %v", err)
			}
		}
	}
}
