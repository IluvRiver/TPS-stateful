package kafka

import (
	"chat_server_golang/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	cfg *config.Config

	producer *kafka.Producer
}

func NewKafka(cfg *config.Config) (*Kafka, error) {
	k := &Kafka{cfg: cfg}

	var err error

	if k.producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URL,
		//부트스트랩 서버는 카프카 url주소가 될것이고
		"client.id": cfg.Kafka.ClientID,
		//클라이언트 아이디는 이 프로듀싱하는 클라이언가 어떤 클라이언트인지 고유한 아이디
		"acks": "all",
		//메시지가 전송되는데 이 고가용성을 위해 복제본을 어디까지 저장하는지
	}); err != nil {
		return nil, err
	} else {
		return k, nil
	}
}

func (k *Kafka) PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	if err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, ch); err != nil {
		return nil, err
	} else {
		return <-ch, nil
	}
}
