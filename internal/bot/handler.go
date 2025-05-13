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

func (h *Handler) Handle(api *tgbotapi.BotAPI, upd tgbotapi.Update) error {
	switch upd.Message.Command() {
	case "distance":
		return h.handleDistance(api, upd)
	default:
		return fmt.Errorf("unknown command: %s", upd.Message.Command())
	}
}

func (h *Handler) handleDistance(api *tgbotapi.BotAPI, upd tgbotapi.Update) error {
	dist, err := data.ReadDistance(h.dataPath)

	var msg tgbotapi.MessageConfig
	if err != nil {
		msg = tgbotapi.NewMessage(
			upd.Message.Chat.ID,
			"Internal error",
		)
	} else {
		text := fmt.Sprintf("Distance: %d", dist)
		msg = tgbotapi.NewMessage(
			upd.Message.Chat.ID,
			text,
		)
	}
	msg.ReplyToMessageID = upd.Message.MessageID
	_, err = api.Send(msg)
	return err
}
