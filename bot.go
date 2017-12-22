package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Syfaro/finch"
	_ "github.com/Syfaro/finch/commands/cancel"
	_ "github.com/Syfaro/finch/commands/help"

	_ "github.com/Syfaro/selectionsbot/commands/manage"
	_ "github.com/Syfaro/selectionsbot/commands/start"
	"github.com/Syfaro/selectionsbot/database"
)

func main() {
	db := sqlx.MustOpen("sqlite3", os.Getenv("DATABASE_PATH"))
	database.DB = db

	db.MustExec(`
		create table if not exists user (
			id integer primary key,
			telegram_id integer not null,
			name text
		);

		create table if not exists selection (
			id integer primary key,
			user_id integer not null,
			chat_id integer not null,
			title text,
			active integer
		);

		create table if not exists selection_item (
			id integer primary key,
			selection_id integer not null,
			item text not null
		);

		create table if not exists selection_vote (
			id integer primary key,
			user_id integer not null,
			selection_id integer not null,
			selection_item_id integer not null
		);
	`)

	token := os.Getenv("TELEGRAM_APITOKEN")

	if token == "" {
		panic("No Telegram token was specified!")
	}

	f := finch.NewFinch(token)

	f.API.Debug = os.Getenv("DEBUG") == "true"

	f.Start()
}
