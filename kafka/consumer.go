package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	topic    string
}

func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	c, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		topic:    topic,
	}, nil
}

func (c *Consumer) Consume(handler func([]byte)) error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return err
	}

	for _, p := range partitions {
		pc, err := c.consumer.ConsumePartition(
			c.topic,
			p,
			sarama.OffsetNewest,
		)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				handler(msg.Value)
			}
		}(pc)
	}

	select {} // block forever
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func PrintJSONConsumer(data []byte) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		log.Println("Invalid JSON:", err)
		return
	}

	out, _ := json.MarshalIndent(v, "", "  ")
	log.Println("Consumed message:\n", string(out))
}
