package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mauzec/ibot-things/internal/data"
)

type Handler struct {
	dataPath string
}

func NewHandler(dataPath string) *Handler {
	return &Handler{
		dataPath: dataPath,
	}
}

func (h *Handler) Handle(b *Bot, upd *tgbotapi.Update) error {
	if !upd.Message.IsCommand() {
		err := BotErrNotCommand
		sendErr := b.SendMessage(upd, "Bot does not support noncommand messages yet", true)
		if sendErr != nil {
			err = fmt.Errorf("[Handle] (1) %v\n(2) %v", err, sendErr)
		}
		return err
	}

	switch upd.Message.Command() {
	case "distance":
		return h.handleDistance(b, upd)
	default:
		err := BotErrNotFoundCommand
		sendErr := b.SendMessage(upd, "Not found this command", true)
		if sendErr != nil {
			err = fmt.Errorf("[Handle] (1) %v\n(2) %v", err, sendErr)
		}
		return err
	}
}

func (h *Handler) handleDistance(b *Bot, upd *tgbotapi.Update) error {
	dist, err := data.ReadDistance(h.dataPath)

	var sendErr error
	if err != nil {
		sendErr = b.SendMessage(upd, "Internal error. Pls try again or later", false)
	} else {
		sendErr = b.SendMessage(upd, fmt.Sprintf("Distance: %d", dist), true)
	}
	if sendErr != nil {
		err = fmt.Errorf("[handleDistance] (1) %v\n(2) %v", err, sendErr)
	}
	return err
}
