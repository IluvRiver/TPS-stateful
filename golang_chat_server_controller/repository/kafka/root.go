package kafka

import (
	"golang_chat_server_controller/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	cfg *config.Config

	consumer *kafka.Consumer
}

func NewKafka(cfg *config.Config) (*Kafka, error) {
	k := &Kafka{cfg: cfg}

	var err error

	if k.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.URL,
		"group.id":          cfg.Kafka.GroupID,
		"auto.offset.reset": "latest", //offset은 어디까지 읽었느냐에 대한 값 latest는 최근
	}); err != nil {
		return nil, err
	} else {
		return k, nil
	}
}
func (k *Kafka) Pool(timeoutMs int) kafka.Event {
	return k.consumer.Poll(timeoutMs)
}

// 토픽이라는건 이 컨슈머는 어떠한 키 값에 들어오는 이벤트를 subscribe할거다
func (k *Kafka) RegisterSubTopic(topic string) error {
	if err := k.consumer.Subscribe(topic, nil); err != nil {
		return err
	}
	return nil
}
