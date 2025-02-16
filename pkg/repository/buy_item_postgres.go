package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type BuyItemPostgres struct {
	db *sql.DB
}

func NewBuyItemPostgres(db *sql.DB) *BuyItemPostgres {
	return &BuyItemPostgres{db: db}
}

func (r *BuyItemPostgres) Buy(userId int, item string) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, nil
	}

	var coins int
	createCoinsQuery := fmt.Sprintf("SELECT coins FROM %s WHERE id = $1", userTable)
	row := tx.QueryRow(createCoinsQuery, userId)
	if err := row.Scan(&coins); err != nil {
		tx.Rollback()
		return 0, err
	}

	var price, itemId int
	createPriceQuery := fmt.Sprintf("SELECT id, price from %s where name = $1", itemsTable)
	row = tx.QueryRow(createPriceQuery, item)
	if err := row.Scan(&itemId, &price); err != nil {
		tx.Rollback()
		return 0, err
	}

	if price > coins {
		tx.Rollback()
		return 0, errors.New("not enough coins")
	}

	var id int
	createSaleQuery := fmt.Sprintf("INSERT INTO %s (user_id, item_id) VALUES ($1, $2) RETURNING id", salesTable)
	row = tx.QueryRow(createSaleQuery, userId, itemId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserUpdateQuery := fmt.Sprintf("UPDATE %s SET coins = coins - $1 WHERE id = $2", userTable)
	_, err = tx.Exec(createUserUpdateQuery, price, userId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
