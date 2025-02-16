package service

import "avito_go/repository"

type SendCoinService struct {
	repository repository.SendCoin
}

func NewSendCoinService(repository repository.SendCoin) *SendCoinService {
	return &SendCoinService{repository: repository}
}

func (s *SendCoinService) Send(userId int, toUser string, amount int) (int, error) {
	return s.repository.Send(userId, toUser, amount)
}
