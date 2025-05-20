package bot

import "errors"

var (
	BotErrBadMsg          = errors.New("provided message is not a command")
	BotErrNotFoundCommand = errors.New("not found this command")
)
