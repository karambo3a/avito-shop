package service

import (
	"avito_go/repository"
	"avito_go/responses"
)

type InformationService struct {
	repository repository.Information
}

func NewInformationService(repository repository.Information) *InformationService {
	return &InformationService{repository: repository}
}

func (s *InformationService) GetInformation(userId int) ([]responses.Transaction, []responses.Transaction, []responses.InventoryItem, int, error) {
	sent, err := s.GetSent(userId)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	recieved, err := s.GetRecieved(userId)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	inventory, err := s.GetInventory(userId)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	return sent, recieved, inventory, 0, nil
}

func (s *InformationService) GetSent(userId int) ([]responses.Transaction, error) {
	return s.repository.GetSent(userId)
}

func (s *InformationService) GetRecieved(userId int) ([]responses.Transaction, error) {
	return s.repository.GetRecieved(userId)
}

func (s *InformationService) GetInventory(userId int) ([]responses.InventoryItem, error) {
	return s.repository.GetInventory(userId)
}

func (s *InformationService) GetCoins(userId int) (int, error) {
	return s.repository.GetCoins(userId)
}
