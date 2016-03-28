package commands

import (
	"github.com/syfaro/finch"
	"gopkg.in/telegram-bot-api.v3"
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

func (cmd start) ShouldExecute(message tgbotapi.Message) bool {
	return finch.SimpleCommand("start", message.Text)
}

func (cmd start) Execute(message tgbotapi.Message) error {
	return cmd.QuickReply(message, usage)
}

func (cmd start) Help() finch.Help {
	return finch.Help{
		Name: "Start",
	}
}
