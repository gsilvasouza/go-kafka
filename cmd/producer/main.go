package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	deliveryChan := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("mensagem", "teste", producer, []byte("transferencia"), deliveryChan)
	DeliveryReport(deliveryChan)

	/* e := <-deliveryChan       //Espera o retorno da mensagem
	msg := e.(*kafka.Message) //Pengando a propriedade msg do evento
	if msg.TopicPartition.Error != nil {
		fmt.Println("Erro ao enviar")
	} else {
		fmt.Println("Mensagem enviada", msg.TopicPartition)
	}
	producer.Flush(1000) */
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "go-kafka-kafka-1:9092",
		"delivery.timeout.ms": "0",
		"acks":                "1",
		"enable.idempotence":  "false", //true -> acks = all
	}
	p, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
	}
	return p
}

func Publish(msg string, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	message := &kafka.Message{
		Value:          []byte(msg),
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            key,
	}

	err :=
		producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}

	return nil
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Erro ao enviar")
			} else {
				fmt.Println("Mensagem enviada", ev.TopicPartition)
			}
		}
	}
}
