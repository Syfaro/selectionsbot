package commands

import (
	"database/sql"

	"github.com/Syfaro/finch"
	"github.com/Syfaro/selectionsbot/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	finch.RegisterCommand(&endSelection{})
}

type endSelection struct {
	finch.CommandBase
}

func (cmd endSelection) ShouldExecute(message tgbotapi.Message) bool {
	return finch.SimpleCommand("end", message.Text)
}

func (cmd endSelection) Execute(message tgbotapi.Message) error {
	var selection database.Selection
	err := database.DB.Get(&selection, `
		select
			*
		from
			selection
		where
			active = 1 and
			chat_id = $1
	`, message.Chat.ID)
	if err == sql.ErrNoRows {
		return cmd.QuickReply(message, "No selection was currently active")
	} else if err != nil {
		return err
	}

	u, err := database.LoadUser(message.From.ID)
	if err == sql.ErrNoRows {
		return cmd.QuickReply(message, "Your user doesn't seem to exist!")
	} else if selection.UserID != u.ID {
		return cmd.QuickReply(message, "You did not create the active selection in this channel")
	}

	_, err = database.DB.Exec(`
		update selection
			set
				active = 0
			where
				chat_id = $1
	`, message.Chat.ID)
	if err != nil {
		return err
	}

	return cmd.QuickReply(message, "Current selection was ended.")
}

func (cmd endSelection) Help() finch.Help {
	return finch.Help{
		Name:        "End",
		Description: "Ends selection",
		Example:     "/end@@",
		Botfather: [][]string{
			[]string{"end", "Ends the current selection"},
		},
	}
}
