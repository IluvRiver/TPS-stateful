package service

import (
	"golang_chat_server_controller/repository"
	"golang_chat_server_controller/types/table"
)

type Service struct {
	repository *repository.Repository

	AvgServerList map[string]bool
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository, AvgServerList: make(map[string]bool)}

	s.setServerInfo()

	return s
}
func (s *Service) GetAvgServerList() []string {
	var res []string
	for ip, available := range s.AvgServerList {
		if available {
			res = append(res, ip)
		}
	}
	return res
}
func (s *Service) setServerInfo() {
	if serverList, err := s.GetAbailableServerList(); err != nil {
		panic(err)
	} else {
		for _, server := range serverList {
			s.AvgServerList[server.IP] = true
		}
	}
}
func (s *Service) GetAbailableServerList() ([]*table.ServerInfo, error) {
	return s.repository.GetAbailableServerList()
}
