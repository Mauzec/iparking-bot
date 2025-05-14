package bot

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mauzec/ibot-things/config"
	"github.com/mauzec/ibot-things/pkg/logger"
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
		logger.Logger().Infof(
			"[BOT] bot timeout is not set or <= 0, using default value: %d",
			cfg.BotTimeout,
		)
	}

	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	logger.Logger().Info("[BOT] USING DEBUG MODE")
	api.Debug = true

	h := NewHandler(cfg.DataPath)

	return &Bot{
		api:           api,
		handler:       h,
		updateTimeout: cfg.BotTimeout,
	}, nil
}

func (b *Bot) Run() {
	// flush all updates; it's good?
	_, err := b.api.Request(tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: true,
	})
	if err != nil {
		logger.Logger().Errorf("[RUN] drop pending updates error: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.updateTimeout

	updates := b.api.GetUpdatesChan(u)
	for upd := range updates {
		if upd.Message != nil {
			if err := b.handler.Handle(b, &upd); err != nil {
				logger.Logger().Debugf("[RUN] %v", err)
			}
		}
	}
}

func (b *Bot) SendMessage(upd *tgbotapi.Update, text string, isReplyToMessage bool) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
	if isReplyToMessage {
		msg.ReplyToMessageID = upd.Message.MessageID
	}
	_, err := b.api.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
