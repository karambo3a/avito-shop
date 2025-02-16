package repository

import (
	"avito_go/pkg/responses"
	"avito_go/pkg/shop"
	"database/sql"
)

type Authorization interface {
	GetUser(username, password string) (shop.User, error)
}

type SendCoin interface {
	Send(userId int, toUser string, amount int) (int, error)
}

type BuyItem interface {
	Buy(userId int, item string) (int, error)
}

type Information interface {
	GetSent(userId int) ([]responses.Transaction, error)
	GetRecieved(userId int) ([]responses.Transaction, error)
	GetInventory(userId int) ([]responses.InventoryItem, error)
	GetCoins(userId int) (int, error)
}

type Repository struct {
	Authorization
	SendCoin
	BuyItem
	Information
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		SendCoin:      NewSendCoinPostgres(db),
		BuyItem:       NewBuyItemPostgres(db),
		Information:   NewInformationPostgres(db),
	}
}
