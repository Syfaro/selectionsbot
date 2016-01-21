package commands

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/syfaro/finch"
	"github.com/syfaro/selectionsbot/database"
)

func init() {
	finch.RegisterCommand(&endSelection{})
}

type endSelection struct {
	finch.CommandBase
}

func (cmd endSelection) ShouldExecute(update tgbotapi.Update) bool {
	return finch.SimpleCommand("end", update.Message.Text)
}

func (cmd endSelection) Execute(update tgbotapi.Update) error {
	var selection database.Selection
	err := database.DB.Get(&selection, `
		select
			*
		from
			selection
		where
			active = 1 and
			chat_id = $1
	`, update.Message.Chat.ID)
	if err != nil {
		return err
	}

	u := database.User{}
	err = u.Load(update.Message.From.ID)
	if err == sql.ErrNoRows || selection.UserID != u.ID {
		return cmd.QuickReply(update.Message,
			"You did not create the active selection in this channel!")
	}

	_, err = database.DB.Exec(`
		update selection
			set
				active = 0
			where
				chat_id = $1
	`, update.Message.Chat.ID)
	if err != nil {
		return err
	}

	return cmd.QuickReply(update.Message, "Current selection was ended.")
}

func (cmd endSelection) Help() finch.Help {
	return finch.Help{
		Name:        "End",
		Description: "Ends selection",
		Example:     "/end",
		Botfather: [][]string{
			[]string{"end", "Ends the current selection"},
		},
	}
}
