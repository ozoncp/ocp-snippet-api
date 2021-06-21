package producer

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
)

type Producer interface {
	SendMessage(message SnippetEvent)
	Close()
}

type producer struct {
	Producer  // нужна ли эта строчка?..
	kafkaProd sarama.SyncProducer
	topic     string
}

type SnippetEventType = int

const (
	Created SnippetEventType = iota
	Updated
	Deleted
)

type SnippetEvent struct {
	Type SnippetEventType
	Body map[string]interface{}
}

var brokers = []string{"127.0.0.1:9094"}

func NewProducer(topic string) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	kafkaProd, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		log.Printf("err: %v", err)
		return nil, err
	}

	return &producer{
		kafkaProd: kafkaProd,
		topic:     topic,
	}, err
}

func (prod *producer) SendMessage(message SnippetEvent) {
	if _, _, err := prod.kafkaProd.SendMessage(prepareMessage(prod.topic, message)); err != nil {
		log.Printf("err: %v", err)
	}
}

func (prod *producer) Close() {
	prod.kafkaProd.Close()
}

func prepareMessage(topic string, message SnippetEvent) *sarama.ProducerMessage {
	dataBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Cannot create json from message: %v", err)
		return nil
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(dataBytes),
	}
	return msg
}
