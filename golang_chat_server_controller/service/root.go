package service

import (
	"encoding/json"
	"golang_chat_server_controller/repository"
	"golang_chat_server_controller/types/table"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Service struct {
	repository *repository.Repository

	AvgServerList map[string]bool
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository, AvgServerList: make(map[string]bool)}

	s.setServerInfo()

	if err := s.repository.Kafka.RegisterSubTopic("chat"); err != nil {
		panic(err)
	} else {
		go s.loopSubKafka()
	}
	return s
}
func (s *Service) loopSubKafka() {
	for {
		ev := s.repository.Kafka.Pool(100)

		switch event := ev.(type) {
		case *kafka.Message:
			type ServerInfoEvent struct {
				IP     string
				Status bool
			}

			var decoder ServerInfoEvent

			if err := json.Unmarshal(event.Value, &decoder); err != nil {
				log.Println("Failed To Decode Event", event.Value)
			} else {
				s.AvgServerList[decoder.IP] = decoder.Status
				log.Println("Success To Set ServerList", decoder.IP, decoder.Status)
			}

		case *kafka.Error:
			log.Println("Failed To Pooling Event", event.Error())
		}
	}
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
