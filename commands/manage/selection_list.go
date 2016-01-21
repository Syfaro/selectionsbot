package commands

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/syfaro/finch"
	"github.com/syfaro/selectionsbot/database"
	"strings"
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

	var items []database.SelectionItem
	err = database.DB.Select(&items, `
		select
			*
		from
			selection_item
		where
			selection_id = $1
	`, selection.ID)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}

	for _, item := range items {
		var votes []database.SelectionVote
		err = database.DB.Select(&votes, `
			select
				*
			from
				selection_vote
			where
				selection_id = $1 and
				selection_item_id = $2
		`, selection.ID, item.ID)
		if err != nil {
			return err
		}

		var users []string

		for _, vote := range votes {
			var user database.User
			err = database.DB.Get(&user, `
				select
					*
				from
					user
				where
					id = $1
			`, vote.UserID)
			if err != nil {
				return err
			}

			users = append(users, user.Name)
		}

		b.WriteString(item.Item)
		b.WriteString(" - ")
		if len(users) == 0 {
			b.WriteString("None")
		} else {
			b.WriteString(strings.Join(users, ", "))
		}
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
