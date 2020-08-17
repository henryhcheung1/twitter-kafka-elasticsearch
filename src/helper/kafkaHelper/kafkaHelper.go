package kafkaHelper

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

type KafkaProperties struct {
	kafkaBrokerUrl     string
	kafkaVerbose       string
	kafkaTopic         string
	kafkaConsumerGroup string
	kafkaClientId      string
}

func Configure(kafkaBrokerUrls []string, clientId string, topic string) (w *kafka.Writer, err error) {

	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}

	config := kafka.WriterConfig{
		Brokers:          kafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}

	w = kafka.NewWriter(config)

	return w, nil
}

func Push(writer *kafka.Writer, parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return writer.WriteMessages(parent, message)
}
