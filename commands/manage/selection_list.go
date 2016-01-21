package commands

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/syfaro/finch"
	"github.com/syfaro/selectionsbot/database"
)

func init() {
	finch.RegisterCommand(&selectionList{})
}

type selectionList struct {
	finch.CommandBase
}

func (cmd selectionList) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("list", update.Message.Text)
}

func (cmd selectionList) Execute(update tgbotapi.Update) error {
	var selection database.Selection
	err := database.DB.Get(&selection, `
		select
			*
		from
			selection
		where
			chat_id = $1
		order by
			id desc
	`, update.Message.Chat.ID)
	if err != nil {
		return err
	}

	var list []database.SelectionList
	err = database.DB.Select(&list, `
		select
			user.name as user,
			selection_item.item as item
		from
			selection_vote
		inner join
			user on
				selection_vote.user_id = user.id
		inner join
			selection_item on
				selection_item.id = selection_vote.selection_item_id
		where
			selection_vote.selection_id = $1
	`, selection.ID)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}

	for _, item := range list {
		b.WriteString(item.User)
		b.WriteString(" - ")
		b.WriteString(item.Item)
		b.WriteString("\n")
	}

	return cmd.QuickReply(update.Message, b.String())
}

func (cmd selectionList) Help() finch.Help {
	return finch.Help{
		Name:        "List",
		Description: "List selected items",
		Example:     "/list@@",
		Botfather: [][]string{
			[]string{"list", "List who has selected what item"},
		},
	}
}
