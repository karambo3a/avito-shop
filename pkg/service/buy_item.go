package service

import "avito_go/pkg/repository"

type BuyItemService struct {
	repository repository.BuyItem
}

func NewBuyItemService(repository repository.BuyItem) *BuyItemService {
	return &BuyItemService{repository: repository}
}

func (s *BuyItemService) Buy(userId int, item string) (int, error) {
	return s.repository.Buy(userId, item)
}
