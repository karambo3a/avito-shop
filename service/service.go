package service

import (
	"avito_go/repository"
	"avito_go/responses"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, int, error)
}

type SendCoin interface {
	Send(userId int, toUser string, amount int) (int, error)
}

type BuyItem interface {
	Buy(userId int, item string) (int, error)
}

type Information interface {
	GetInformation(userId int) ([]responses.Transaction, []responses.Transaction, []responses.InventoryItem, int, error)
}

type Service struct {
	Authorization
	SendCoin
	BuyItem
	Information
}


func NewService(repository *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repository.Authorization),
		SendCoin:      NewSendCoinService(repository.SendCoin),
		BuyItem:       NewBuyItemService(repository.BuyItem),
		Information:   NewInformationService(repository.Information),
	}
}
