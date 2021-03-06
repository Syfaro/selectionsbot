package commands

import (
	"database/sql"

	"github.com/Syfaro/finch"
	"github.com/Syfaro/selectionsbot/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	finch.RegisterCommand(&selectItem{})
}

type selectItem struct {
	finch.CommandBase
}

func (cmd selectItem) ShouldExecute(message tgbotapi.Message) bool {
	return finch.SimpleCommand("select", message.Text)
}

func (cmd selectItem) Execute(message tgbotapi.Message) error {
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
	if err == sql.ErrNoRows {
		return cmd.QuickReply(message, "There are no active selections")
	} else if err != nil {
		return err
	}

	var items []database.SelectionItem
	err = database.DB.Select(&items, `
		select
			*
		from
			selection_item
		where
			selection_id = $1
	`, selection.ID)
	if err == sql.ErrNoRows {
		return cmd.QuickReply(message, "No items were added to the current selection")
	} else if err != nil {
		return err
	}

	var itemList [][]tgbotapi.KeyboardButton

	for _, item := range items {
		itemList = append(itemList, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(item.Item)))
	}

	msg := tgbotapi.NewMessage(message.Chat.ID,
		"Select your item")

	if selection.Title.Valid && selection.Title.String != "" {
		msg.Text = msg.Text + " for: " + selection.Title.String
	}

	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        itemList,
		Selective:       true,
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	cmd.SetWaiting(message.From.ID)

	return cmd.SendMessage(msg)
}

func (cmd selectItem) ExecuteWaiting(message tgbotapi.Message) error {
	cmd.ReleaseWaiting(message.From.ID)

	user, err := database.LoadUser(message.From.ID)
	if err == sql.ErrNoRows {
		if err = user.Init(message.From); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var selection database.Selection
	err = database.DB.Get(&selection, `
		select
			*
		from
			selection
		where
			chat_id = $1 and
			active = 1
	`, message.Chat.ID)
	if err == sql.ErrNoRows {
		return cmd.QuickReply(message, "There are no active selections in this chat currently")
	} else if err != nil {
		return err
	}

	var item database.SelectionItem
	err = database.DB.Get(&item, `
		select
			*
		from
			selection_item
		where
			item = $1 and
			selection_id = $2
	`, message.Text, selection.ID)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(`
		delete from
			selection_vote
		where
			user_id = $1 and
			selection_id = $2
	`, user.ID, selection.ID)
	if err != nil {
		return err
	}

	_, err = database.NewSelectionVote(user.ID, selection.ID, item.ID)
	if err != nil {
		return err
	}

	return cmd.QuickReply(message, "Added selection!")
}

func (cmd selectItem) Help() finch.Help {
	return finch.Help{
		Name:        "Select item",
		Description: "Selects an item",
		Example:     "/select@@",
		Botfather: [][]string{
			[]string{"select", "Select an item"},
		},
	}
}
