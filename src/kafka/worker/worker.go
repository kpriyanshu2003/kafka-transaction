package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"example.com/kafka/src/models"
	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
)

var (
	db       *sql.DB
	consumer sarama.ConsumerGroup
	wg       sync.WaitGroup
)

func init() {
	// Initialize Kafka consumer group
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V3_9_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	var err error
	consumer, err = sarama.NewConsumerGroup([]string{"localhost:9092"}, "wallet-service-group", config)
	if err != nil {
		log.Panicln("Failed to create Kafka consumer group:", err)
	}
	fmt.Println("Connected to Kafka")

	// Initialize PostgreSQL connection
	db, err = sql.Open("postgres", "user=postgres password=docker dbname=kafka sslmode=disable")
	if err != nil {
		log.Panicln("Failed to connect to PostgreSQL:", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	fmt.Println("Connected to PostgreSQL")
}

func main() {
	topic := "wallet.pending"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumer.Consume(ctx, []string{topic}, &MessageHandler{}); err != nil {
				if ctx.Err() != nil {
					return // Context canceled, exit goroutine
				}
				log.Println("Error consuming messages:", err)
			}
		}
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan
	fmt.Println("Interrupt detected. Shutting down...")
	cancel()

	wg.Wait() // Wait for all goroutines to finish

	// Cleanup
	if err := consumer.Close(); err != nil {
		log.Println("Error closing Kafka consumer group:", err)
	}
	if err := db.Close(); err != nil {
		log.Println("Error closing PostgreSQL connection:", err)
	}
}

type MessageHandler struct{}

func (h *MessageHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (h *MessageHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *MessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var txnMessage models.Transaction
		if err := json.Unmarshal(message.Value, &txnMessage); err != nil {
			log.Printf("Failed to parse message: %v\n", err)
			moveMessage("wallet.failed", message.Value)
			continue
		}

		if processTransaction(txnMessage) {
			moveMessage("wallet.done", message.Value)
		} else {
			moveMessage("wallet.failed", message.Value)
		}

		session.MarkMessage(message, "")
	}
	return nil
}

func processTransaction(txn models.Transaction) bool {
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return false
	}

	switch txn.Type {
	case "CREDIT":
		_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE user_id = $2", txn.Amt, txn.SendID)
		if err != nil {
			log.Println("Failed to credit balance:", err)
			tx.Rollback()
			return false
		}

	case "DEBIT":
		_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE user_id = $2 AND balance >= $1", txn.Amt, txn.SendID)
		if err != nil {
			log.Println("Failed to debit balance or insufficient funds:", err)
			tx.Rollback()
			return false
		}

	case "MOVE":
		_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE user_id = $2 AND balance >= $1", txn.Amt, txn.SendID)
		if err != nil {
			log.Println("Failed to subtract balance during move:", err)
			tx.Rollback()
			return false
		}

		_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE user_id = $2", txn.Amt, txn.RecvID)
		if err != nil {
			log.Println("Failed to add balance during move:", err)
			tx.Rollback()
			return false
		}

	default:
		log.Println("Invalid transaction type:", txn.Type)
		tx.Rollback()
		return false
	}

	_, err = tx.Exec("INSERT INTO transactions (send_id, recv_id, amt, status, type) VALUES ($1, $2, $3, $4, $5)",
		txn.SendID, txn.RecvID, txn.Amt, "DONE", txn.Type)
	if err != nil {
		log.Println("Failed to insert transaction:", err)
		tx.Rollback()
		return false
	}

	if err := tx.Commit(); err != nil {
		log.Println("Failed to commit transaction:", err)
		return false
	}

	return true
}

func moveMessage(topic string, message []byte) {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Println("Failed to create Kafka producer:", err)
		return
	}
	defer producer.Close()

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	})
	if err != nil {
		log.Printf("Failed to move message to topic %s: %v\n", topic, err)
	}
}
