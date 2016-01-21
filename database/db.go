package database

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/telegram-bot-api.v2"
)

var DB *sqlx.DB

type User struct {
	ID         int64  `db:"id"`
	TelegramID int    `db:"telegram_id"`
	Name       string `db:"name"`
}

func (u *User) Load(telegramID int) error {
	return DB.Get(u, `
		select
			*
		from
			user
		where
			telegram_id = $1
	`, telegramID)
}

func (u *User) Init(user tgbotapi.User) error {
	u.TelegramID = user.ID
	u.Name = user.String()

	res, err := DB.Exec(`
		insert into user
			(telegram_id, name) values
			($1, $2)
	`, u.TelegramID, u.Name)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id

	return nil
}

type Selection struct {
	ID     int64          `db:"id"`
	UserID int64          `db:"user_id"`
	ChatID int            `db:"chat_id"`
	Title  sql.NullString `db:"title"`
	Active bool           `db:"active"`
}

func NewSelection(userID int64, chatID int) (Selection, error) {
	return NewSelectionWithTitle(userID, chatID, nil)
}

func NewSelectionWithTitle(userID int64, chatID int, title *string) (Selection, error) {
	res, err := DB.Exec(`
		insert into selection
			(user_id, chat_id, active) values
			($1, $2, 1)
	`, userID, chatID)
	if err != nil {
		return Selection{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Selection{}, err
	}

	return Selection{
		ID:     id,
		UserID: userID,
		ChatID: chatID,
		Active: true,
	}, nil
}

type SelectionItem struct {
	ID          int64  `db:"id"`
	SelectionID int64  `db:"selection_id"`
	Item        string `db:"item"`
}

func NewSelectionItem(selectionID int64, item string) (SelectionItem, error) {
	res, err := DB.Exec(`
		insert into selection_item
			(selection_id, item) values
			($1, $2)
	`, selectionID, item)
	if err != nil {
		return SelectionItem{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return SelectionItem{}, err
	}

	return SelectionItem{
		ID:          id,
		SelectionID: selectionID,
		Item:        item,
	}, nil
}

type SelectionVote struct {
	ID              int64 `db:"id"`
	UserID          int64 `db:"user_id"`
	SelectionID     int64 `db:"selection_id"`
	SelectionItemID int64 `db:"selection_item_id"`
}

func NewSelectionVote(userID, selectionID, selectionItemID int64) (SelectionVote, error) {
	res, err := DB.Exec(`
		insert into selection_vote
			(user_id, selection_id, selection_item_id) values
			($1, $2, $3)
	`, userID, selectionID, selectionItemID)
	if err != nil {
		return SelectionVote{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return SelectionVote{}, err
	}

	return SelectionVote{
		ID:              id,
		UserID:          userID,
		SelectionID:     selectionID,
		SelectionItemID: selectionItemID,
	}, nil
}

type SelectionVoteCount struct {
	Count int64  `db:"count"`
	Item  string `db:"item"`
}

type SelectionList struct {
	User string `db:"user"`
	Item string `db:"item"`
}
