package commands

import (
	"bytes"
	"database/sql"
	"strings"

	"github.com/Syfaro/finch"
	"github.com/Syfaro/selectionsbot/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	finch.RegisterCommand(&selectionList{})
}

type selectionList struct {
	finch.CommandBase
}

func (cmd selectionList) ShouldExecute(message tgbotapi.Message) bool {
	return finch.SimpleCommand("list", message.Text)
}

func (cmd selectionList) Execute(message tgbotapi.Message) error {
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
	`, message.Chat.ID)
	if err == sql.ErrNoRows {
		return cmd.QuickReply(message, "There are no active selections in this chat currently")
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
		return cmd.QuickReply(message, "Somehow, the selected items disappeared")
	} else if err != nil {
		return err
	}

	b := bytes.Buffer{}

	if selection.Title.Valid && selection.Title.String != "" {
		b.WriteString("Counts from ")
		b.WriteString(selection.Title.String)
		b.WriteString("\n\n")
	}

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
		if err != sql.ErrNoRows {
			return cmd.QuickReply(message, "No selections have been cast yet")
		} else if err != nil {
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
			if err == sql.ErrNoRows {
				user.Name = "Deleted user"
			} else if err != nil {
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

	return cmd.QuickReply(message, b.String())
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
