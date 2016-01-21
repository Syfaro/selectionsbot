package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syfaro/finch"
	_ "github.com/syfaro/finch/commands/help"
	_ "github.com/syfaro/selectionsbot/commands/manage"
	_ "github.com/syfaro/selectionsbot/commands/start"
	"github.com/syfaro/selectionsbot/database"
	"os"
)

func main() {
	db, err := sqlx.Open("sqlite3", os.Getenv("DATABASE_PATH"))
	if err != nil {
		panic(err)
	}

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

	f := finch.NewFinch(os.Getenv("TELEGRAM_APITOKEN"))

	f.API.Debug = true

	f.Start()
}
