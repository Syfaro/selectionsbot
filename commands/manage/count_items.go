package commands

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/syfaro/finch"
	"github.com/syfaro/selectionsbot/database"
	"strconv"
)

func init() {
	finch.RegisterCommand(&countItems{})
}

type countItems struct {
	finch.CommandBase
}

func (cmd countItems) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("count", update.Message.Text)
}

func (cmd countItems) Execute(update tgbotapi.Update) error {
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
	if err != nil {
		return err
	}

	b := bytes.Buffer{}

	for _, item := range items {
		b.WriteString(item.Item)
		b.WriteString(" - ")
		b.WriteString(strconv.FormatInt(item.Count, 10))
		b.WriteString("\n")
	}

	return cmd.QuickReply(update.Message, b.String())
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
