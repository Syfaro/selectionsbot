package commands

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/syfaro/finch"
	"github.com/syfaro/selectionsbot/database"
	"strconv"
	"strings"
)

func init() {
	finch.RegisterCommand(&createSelection{})
}

type createSelection struct {
	finch.CommandBase
}

func (cmd createSelection) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("create", update.Message.Text)
}

func (cmd createSelection) Execute(update tgbotapi.Update) error {
	var selection database.Selection
	err := database.DB.Get(&selection, `
		select
			*
		from
			selection
		where
			chat_id = $1 and
			active = 1
	`, update.Message.Chat.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err != sql.ErrNoRows {
		return cmd.QuickReply(update.Message,
			"Please end the current selection in this chat first!")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Please enter a list of items seperated by new lines.")

	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
		Selective:  true,
	}

	msg.ReplyToMessageID = update.Message.MessageID

	cmd.SetWaiting(update.Message.From.ID)

	return cmd.SendMessage(msg)
}

func (cmd createSelection) ExecuteKeyboard(update tgbotapi.Update) error {
	cmd.ReleaseWaiting(update.Message.From.ID)

	user := database.User{}
	err := user.Load(update.Message.From.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else if err == sql.ErrNoRows {
		if err = user.Init(update.Message.From); err != nil {
			return err
		}
	}

	items := strings.Split(update.Message.Text, "\n")

	selection, err := database.NewSelection(user.ID, update.Message.Chat.ID)
	if err != nil {
		return err
	}

	num := 0

	for _, item := range items {
		_, err := database.NewSelectionItem(selection.ID, item)
		if err != nil {
			return err
		}

		num++
	}

	return cmd.QuickReply(update.Message, "Added "+strconv.Itoa(num)+" items to selection!")
}

func (cmd createSelection) Help() finch.Help {
	return finch.Help{
		Name:        "Create selection",
		Description: "Create a new selection set",
		Example:     "/create",
		Botfather: [][]string{
			[]string{"create", "Creates a new selection set"},
		},
	}
}
