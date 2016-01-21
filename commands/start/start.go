package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/syfaro/finch"
)

var usage string = `How to use this bot:

Invite it to your channel, then do /create.
After inputting all your items in the reply, users may select an item with /select.

When finished, you can get total counts with /count or who selected what with /list. Run /end to create a new selection.`

func init() {
	finch.RegisterCommand(&start{})
}

type start struct {
	finch.CommandBase
}

func (cmd start) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("start", update.Message.Text)
}

func (cmd start) Execute(update tgbotapi.Update) error {
	return cmd.QuickReply(update.Message, usage)
}
