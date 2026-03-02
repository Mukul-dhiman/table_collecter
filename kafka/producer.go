package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: p,
		topic:    topic,
	}, nil
}

func (p *Producer) Enqueue(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Println("Kafka: message enqueued")
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
