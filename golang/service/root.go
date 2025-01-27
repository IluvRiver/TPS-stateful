package service

import (
	"chat_server_golang/repository"
	"chat_server_golang/types/schema"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository}

	return s
}

func (s *Service) PublishServerStatusEvent(ip string, status bool) {
	//kafka에 이벤트 전송
	type ServerInfoEvent struct {
		IP     string
		Status bool
	}
	e := &ServerInfoEvent{IP: ip, Status: status}
	ch := make(chan kafka.Event)

	if v, err := json.Marshal(e); err != nil {
		log.Println("Failed to Marshal")
	} else if result, err := s.PublishEvent("chat", v, ch); err != nil {
		//Todo send event to Kafka
		log.Println("Failed to Send Event to Kafka", "err", err)
	} else {
		log.Println("Successfully Send Event to Kafka", result)
	}
}
func (s *Service) PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	return s.repository.Kafka.PublishEvent(topic, value, ch)
}

func (s *Service) ServerSet(ip string, available bool) error {
	if err := s.repository.ServerSet(ip, available); err != nil {
		log.Println("Failed To ServerSet", "ip", ip, "available", available)
		return err
	} else {
		return nil
	}
}
func (s *Service) InsertChatting(user, message, roomName string) {
	if err := s.repository.InsertChatting(user, message, roomName); err != nil {
		log.Println("Failed To Chat", "err", err)
	}
}

func (s *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	// TODO 필요하다면, 채팅방에 대한 추가 적인 정보를 불러 올 수 있게!
	if res, err := s.repository.GetChatList(roomName); err != nil {
		log.Println("Failed To Get Chat List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}
func (s *Service) GetChatList(roomName string) ([]*schema.Chat, error) {
	if res, err := s.repository.GetChatList(roomName); err != nil {
		log.Println("Failed To Get chat List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) RoomList() ([]*schema.Room, error) {
	if res, err := s.repository.RoomList(); err != nil {
		log.Println("Failed To Get All Room List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) MakeRoom(name string) error {
	if err := s.repository.MakeRoom(name); err != nil {
		log.Println("Failed To Make New Room", "err", err.Error())
		return err
	} else {
		return nil
	}
}

func (s *Service) Room(name string) (*schema.Room, error) {
	if res, err := s.repository.Room(name); err != nil {
		log.Println("Failed To Get Room ", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}
