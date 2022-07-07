package worker

import (
	"context"
	"encoding/json"
	"github.com/abrbird/orders-tracking/internal/metrics"
	"log"

	"github.com/Shopify/sarama"
	cnfg "github.com/abrbird/orders-tracking/config"
	"github.com/abrbird/orders-tracking/internal/broker/kafka"
	"github.com/abrbird/orders-tracking/internal/models"
	rpstr "github.com/abrbird/orders-tracking/internal/repository"
	srvc "github.com/abrbird/orders-tracking/internal/service"
	"github.com/pkg/errors"
)

type IssueOrderHandler struct {
	producer   sarama.SyncProducer
	repository rpstr.Repository
	service    srvc.Service
	metrics    metrics.Metrics
	config     *cnfg.Config
}

func (i *IssueOrderHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (i *IssueOrderHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (i *IssueOrderHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := context.Background()

		if msg.Topic != i.config.Kafka.IssueOrderTopics.IssueOrder {
			log.Printf(
				"topic names does not match: expected %s, got %s\n",
				i.config.Kafka.IssueOrderTopics.IssueOrder,
				msg.Topic,
			)
			continue
		}

		var issueOrderMessage kafka.IssueOrderMessage
		err := json.Unmarshal(msg.Value, &issueOrderMessage)
		if err != nil {
			i.metrics.Error()
			log.Print("Unmarshall failed: value=%v, err=%v", string(msg.Value), err)
			continue
		}

		log.Printf("consumer %s: <- %s: %v",
			i.config.Application.Name,
			i.config.Kafka.IssueOrderTopics.IssueOrder,
			issueOrderMessage,
		)

		record := i.service.OrderHistory().IssueOrder(ctx, i.repository.OrderHistory(), issueOrderMessage.Order.Id)
		if record.Error != nil {
			i.metrics.Error()
			if errors.Is(record.Error, models.RetryError) {
				err = i.RetryIssueOrder(issueOrderMessage)
				if err != nil {
					i.metrics.KafkaError()
					log.Println(err)
				} else {
					log.Printf(
						"consumer %s: -> %s: %v",
						i.config.Application.Name,
						i.config.Kafka.IssueOrderTopics.IssueOrder,
						issueOrderMessage,
					)
				}
			} else {
				i.metrics.KafkaError()
				log.Println(record.Error)
			}
			continue
		}

		err = i.SendRemoveOrder(issueOrderMessage)
		if err != nil {
			i.metrics.KafkaError()
			log.Println(err)
		} else {
			log.Printf(
				"consumer %s: -> %s: %v",
				i.config.Application.Name,
				i.config.Kafka.IssueOrderTopics.RemoveOrder,
				issueOrderMessage,
			)
		}
	}
	return nil
}

func (i *IssueOrderHandler) RetryIssueOrder(message kafka.IssueOrderMessage) error {
	message.Base.SenderServiceName = i.config.Application.Name
	message.Base.Attempt += 1
	maxAttempts := int64(5)

	if message.Base.Attempt > maxAttempts {
		return models.MaxAttemptsError(nil, maxAttempts)
	}

	part, offs, kerr, err := kafka.SendMessage(i.producer, i.config.Kafka.IssueOrderTopics.IssueOrder, message)
	if err != nil {
		return models.BrokerSendError(err)
	}

	if kerr != nil {
		return models.BrokerSendError(err)
	}

	_ = part
	_ = offs

	return nil
}

func (i *IssueOrderHandler) SendRemoveOrder(message kafka.IssueOrderMessage) error {

	message.Base.SenderServiceName = i.config.Application.Name

	part, offs, kerr, err := kafka.SendMessage(i.producer, i.config.Kafka.IssueOrderTopics.RemoveOrder, message)
	if err != nil {
		return models.BrokerSendError(err)
	}

	if kerr != nil {
		return models.BrokerSendError(err)
	}

	_ = part
	_ = offs

	return nil
}
