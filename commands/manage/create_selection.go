package commands

import (
	"database/sql"
	"github.com/syfaro/finch"
	"github.com/syfaro/selectionsbot/database"
	"gopkg.in/telegram-bot-api.v2"
	"strconv"
	"strings"
)

func init() {
	finch.RegisterCommand(&createSelection{})
}

type createSelection struct {
	finch.CommandBase
}

func (cmd createSelection) ShouldExecute(message tgbotapi.Message) bool {
	return finch.SimpleCommand("create", message.Text)
}

func (cmd createSelection) Execute(message tgbotapi.Message) error {
	var selection database.Selection
	err := database.DB.Get(&selection, `
		select
			*
		from
			selection
		where
			chat_id = $1 and
			active = 1
	`, message.Chat.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err != sql.ErrNoRows {
		return cmd.QuickReply(message,
			"Please end the current selection in this chat first!")
	}

	msg := tgbotapi.NewMessage(message.Chat.ID,
		"Please enter a list of items seperated by new lines.\nAn item that is prefixed with a '!' will become the title.")

	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
		Selective:  true,
	}

	msg.ReplyToMessageID = message.MessageID

	cmd.SetWaiting(message.From.ID)

	return cmd.SendMessage(msg)
}

func (cmd createSelection) ExecuteWaiting(message tgbotapi.Message) error {
	cmd.ReleaseWaiting(message.From.ID)

	user := database.User{}
	err := user.Load(message.From.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err == sql.ErrNoRows {
		if err = user.Init(message.From); err != nil {
			return err
		}
	}

	items := strings.Split(message.Text, "\n")

	selection, err := database.NewSelection(user.ID, message.Chat.ID)
	if err != nil {
		return err
	}

	num := 0

	for _, item := range items {
		if item[0] == '!' {
			_, err = database.DB.Exec(`
				update selection
					set title = $1
					where id = $2
			`, item[1:], selection.ID)
			if err != nil {
				return err
			}

			continue
		}

		_, err := database.NewSelectionItem(selection.ID, item)
		if err != nil {
			return err
		}

		num++
	}

	return cmd.QuickReply(message, "Added "+strconv.Itoa(num)+" items to selection!")
}

func (cmd createSelection) Help() finch.Help {
	return finch.Help{
		Name:        "Create selection",
		Description: "Create a new selection set",
		Example:     "/create@@",
		Botfather: [][]string{
			[]string{"create", "Creates a new selection set"},
		},
	}
}
