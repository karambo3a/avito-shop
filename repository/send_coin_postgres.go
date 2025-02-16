package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type SendCoinPostgres struct {
	db *sql.DB
}

func NewSendCoinPostgres(db *sql.DB) *SendCoinPostgres {
	return &SendCoinPostgres{db: db}
}

func (r *SendCoinPostgres) Send(userId int, toUser string, amount int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var coins int
	createUserCoinsQuery := fmt.Sprintf("SELECT coins FROM %s WHERE id = $1", userTable)
	row := tx.QueryRow(createUserCoinsQuery, userId)

	if err := row.Scan(&coins); err != nil {
		tx.Rollback()
		return 0, err
	}

	if coins < amount {
		tx.Rollback()
		return 0, errors.New("not enough coins")
	}

	createFromUserQuery := fmt.Sprintf("UPDATE %s SET coins = coins - $1 WHERE id = $2", userTable)
	_, err = tx.Exec(createFromUserQuery, amount, userId)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var toUserId int
	createToUserQuery := fmt.Sprintf("UPDATE %s SET coins = coins + $1 WHERE username = $2 RETURNING id", userTable)
	row = tx.QueryRow(createToUserQuery, amount, toUser)

	if err := row.Scan(&toUserId); err != nil {
		tx.Rollback()
		return 0, err
	}

	log.Default().Println(userId)

	var id int
	createTransactionQuery := fmt.Sprintf("INSERT INTO %s (from_user_id, to_user_id, amount) VALUES ($1, $2, $3) RETURNING id", transactionsTable)
	row = tx.QueryRow(createTransactionQuery, userId, toUserId, amount)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
