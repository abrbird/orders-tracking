package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

func NewSyncProducer(brokersAddresses []string, brokerCfg *sarama.Config) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(brokersAddresses, brokerCfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func SendMessage(producer sarama.SyncProducer, topic string, message IssueOrderMessage) (int32, int64, error, error) {

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return -1, -1, nil, err
	}

	par, off, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%v", message.Base.StartedAt)),
		Value: sarama.ByteEncoder(messageBytes),
	})

	return par, off, err, nil
}
