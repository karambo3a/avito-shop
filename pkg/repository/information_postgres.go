package repository

import (
	"avito_go/pkg/responses"
	"database/sql"
	"fmt"
)

type InformationPostgres struct {
	db *sql.DB
}

func NewInformationPostgres(db *sql.DB) *InformationPostgres {
	return &InformationPostgres{db: db}
}

func (r *InformationPostgres) GetSent(userId int) ([]responses.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	transactions := make([]responses.Transaction, 0)
	createSentQuery := fmt.Sprintf("SELECT username, amount FROM %s AS t JOIN %s AS u ON t.to_user_id = u.id WHERE from_user_id = $1;", transactionsTable, userTable)
	rows, err := tx.Query(createSentQuery, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction responses.SentTransaction
		if err := rows.Scan(&transaction.ToUser, &transaction.Amount); err != nil {
			tx.Rollback()
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return transactions, tx.Commit()
}

func (r *InformationPostgres) GetRecieved(userId int) ([]responses.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	transactions := make([]responses.Transaction, 0)
	createRecievedQuery := fmt.Sprintf("SELECT username, amount FROM %s AS t JOIN %s AS u ON t.from_user_id = u.id WHERE to_user_id = $1;", transactionsTable, userTable)
	rows, err := tx.Query(createRecievedQuery, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction responses.RecievedTransaction
		if err := rows.Scan(&transaction.FromUser, &transaction.Amount); err != nil {
			tx.Rollback()
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return transactions, tx.Commit()
}

func (r *InformationPostgres) GetInventory(userId int) ([]responses.InventoryItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	inventory := make([]responses.InventoryItem, 0)
	createInventoryQuery := fmt.Sprintf("SELECT i.name, count(i.name) FROM %s AS s JOIN %s AS i ON s.item_id = i.id WHERE user_id = $1 GROUP BY i.name", salesTable, itemsTable)
	rows, err := tx.Query(createInventoryQuery, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var inventoryItem responses.InventoryItem
		if err := rows.Scan(&inventoryItem.Type, &inventoryItem.Quantity); err != nil {
			tx.Rollback()
			return nil, err
		}
		inventory = append(inventory, inventoryItem)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return inventory, tx.Commit()
}

func (r *InformationPostgres) GetCoins(userId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var coins int
	createCoinsQuery := fmt.Sprintf("SELECT coins FROM %s WHERE user_id = $1 ", userTable)
	row := tx.QueryRow(createCoinsQuery, userId)

	if err := row.Scan(&coins); err != nil {
		tx.Rollback()
		return 0, err
	}

	return coins, nil
}
