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
	text := upd.Message.Text
	cmd := upd.Message.Command()

	switch {
	case cmd == "start":
		kb := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Distance"),
				tgbotapi.NewKeyboardButton("Check"),
			),
		)
		kb.ResizeKeyboard = true

		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")
		msg.ReplyMarkup = kb
		_, err := b.api.Send(msg)
		if err != nil {
			err = fmt.Errorf("[Handle] %v", err)
		}
		return err

	case cmd == "distance" || text == "Distance":
		return h.handleDistance(b, upd)
	case cmd == "check" || text == "Check":
		return h.handleCheck(b, upd)
	}

	if upd.Message.IsCommand() {
		err := BotErrNotFoundCommand
		sendErr := b.SendMessage(upd, "Not found this command", true)
		if sendErr != nil {
			err = fmt.Errorf("[Handle] (1) %v\n(2) %v", err, sendErr)
		}
		return err
	}

	err := BotErrBadMsg
	sendErr := b.SendMessage(upd, "Bot does not recognize your message", true)
	if sendErr != nil {
		err = fmt.Errorf("[Handle] (1) %v\n(2) %v", err, sendErr)
	}
	return err

}
func (h *Handler) handleCheck(b *Bot, upd *tgbotapi.Update) error {
	dist, err := data.ReadDistance(h.dataPath)

	var sendErr error
	if err != nil {
		sendErr = b.SendMessage(upd, "Internal error. Pls try again or later", false)
	} else {
		text := "P1 is "
		if dist <= 20 {
			text += "not available"
		} else {
			text += "available"
		}
		sendErr = b.SendMessage(upd, text, true)
	}
	if sendErr != nil {
		err = fmt.Errorf("[handleDistance] (1) %v\n(2) %v", err, sendErr)
	}
	return err
}
func (h *Handler) handleDistance(b *Bot, upd *tgbotapi.Update) error {
	dist, err := data.ReadDistance(h.dataPath)

	var sendErr error
	if err != nil {
		sendErr = b.SendMessage(upd, "Internal error. Pls try again or later", false)
	} else {
		sendErr = b.SendMessage(upd, fmt.Sprintf("Distance: %v", dist), true)
	}
	if sendErr != nil {
		err = fmt.Errorf("[handleDistance] (1) %v\n(2) %v", err, sendErr)
	}
	return err
}
