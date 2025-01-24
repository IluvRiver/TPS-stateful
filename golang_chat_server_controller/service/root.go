package service

import (
	"golang_chat_server_controller/repository"
	"golang_chat_server_controller/types/table"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository}

	return s
}
func (s *Service) GetAbailableServerList() ([]*table.ServerInfo, error) {
	return s.repository.GetAbailableServerList()
}
