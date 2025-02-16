package repository

import (
	"avito_go/shop"
	"database/sql"
	"fmt"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(username, password string) (shop.User, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password, coins) values ($1, $2, $3) RETURNING id", userTable)

	row := r.db.QueryRow(query, username, password, 1000)
	if err := row.Scan(&id); err != nil {
		return shop.User{}, err
	}

	return shop.User{Id: id, Username: username, Password: password, Coins: 1000}, nil
}

func (r *AuthPostgres) GetUser(username, password string) (shop.User, error) {
	var user shop.User
	query := fmt.Sprintf("SELECT id, coins FROM %s WHERE username=$1 AND password=$2", userTable)

	row := r.db.QueryRow(query, username, password)

	if err := row.Scan(&user.Id, &user.Coins); err != nil {
		if err == sql.ErrNoRows {
			return r.CreateUser(username, password)
		}
		return shop.User{}, err
	}

	return user, nil
}
