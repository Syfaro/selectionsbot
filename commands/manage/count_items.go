package commands

import (
	"bytes"
	"database/sql"
	"strconv"

	"github.com/Syfaro/finch"
	"github.com/Syfaro/selectionsbot/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	finch.RegisterCommand(&countItems{})
}

type countItems struct {
	finch.CommandBase
}

func (cmd countItems) ShouldExecute(message tgbotapi.Message) bool {
	return finch.SimpleCommand("count", message.Text)
}

func (cmd countItems) Execute(message tgbotapi.Message) error {
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
		return cmd.QuickReply(message, "Wow, no selections have been made for this channel yet! Try creating one with /create@@")
	} else if err != nil {
		return err
	}

	var items []database.SelectionVoteCount
	err = database.DB.Select(&items, `
		select
			count(selection_vote.id) as count,
			selection_item.item
		from
			selection_vote
		inner join
			selection_item on
				selection_vote.selection_item_id = selection_item.id
		where
			selection_vote.selection_id = $1
		group by
			selection_item_id
	`, selection.ID)
	if err == sql.ErrNoRows {
		return cmd.QuickReply(message, "No selections have been cast yet")
	} else if err != nil {
		return err
	}

	if len(items) == 0 {
		return cmd.QuickReply(message, "No selections have been cast yet")
	}

	b := bytes.Buffer{}

	if selection.Title.Valid && selection.Title.String != "" {
		b.WriteString("Counts from ")
		b.WriteString(selection.Title.String)
		b.WriteString("\n\n")
	}

	for _, item := range items {
		b.WriteString(item.Item)
		b.WriteString(" - ")
		b.WriteString(strconv.FormatInt(item.Count, 10))
		b.WriteString("\n")
	}

	return cmd.QuickReply(message, b.String())
}

func (cmd countItems) Help() finch.Help {
	return finch.Help{
		Name:        "Counts",
		Description: "Total counts of selections",
		Example:     "/count@@",
		Botfather: [][]string{
			[]string{"count", "Count total selections on each item"},
		},
	}
}
