package producer

import (
	"fmt"

	"example.com/kafka/src/config"
	"github.com/IBM/sarama"
)

func PushTransactionToQueue(topic string, jsonMessage []byte) error {
	producer := config.Producer.Producer

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonMessage),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
