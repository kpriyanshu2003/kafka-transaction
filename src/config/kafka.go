package config

import (
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

type ProducerInstance struct {
	Producer sarama.SyncProducer
}

var Producer ProducerInstance

func SetupKafkaProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(GetKafkaBrokers(), config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	log.Println("Kafka producer connected")
	Producer = ProducerInstance{Producer: conn}
}

func GetKafkaBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		log.Fatal("KAFKA_BROKERS is not found in the environment")
	}
	brokersUrl := strings.Split(brokers, ",")
	return brokersUrl
}
