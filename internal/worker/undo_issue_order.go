package worker

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	cnfg "github.com/abrbird/orders-tracking/config"
	"github.com/abrbird/orders-tracking/internal/broker/kafka"
	"github.com/abrbird/orders-tracking/internal/metrics"
	"github.com/abrbird/orders-tracking/internal/models"
	"log"

	rpstr "github.com/abrbird/orders-tracking/internal/repository"
	srvc "github.com/abrbird/orders-tracking/internal/service"
)

type UndoIssueOrderHandler struct {
	producer   sarama.SyncProducer
	repository rpstr.Repository
	service    srvc.Service
	metrics    metrics.Metrics
	config     *cnfg.Config
}

func (u *UndoIssueOrderHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (u *UndoIssueOrderHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (u *UndoIssueOrderHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := context.Background()

		if msg.Topic != u.config.Kafka.IssueOrderTopics.UndoIssueOrder {
			log.Printf(
				"topic names does not match: expected - %s, got %s\n",
				u.config.Kafka.IssueOrderTopics.UndoIssueOrder,
				msg.Topic,
			)
			continue
		}

		var issueOrderMessage kafka.IssueOrderMessage
		err := json.Unmarshal(msg.Value, &issueOrderMessage)
		if err != nil {
			u.metrics.Error()
			log.Print("Unmarshall failed: value=%v, err=%v", string(msg.Value), err)
			continue
		}

		log.Printf("consumer %s: <- %s: %v",
			u.config.Application.Name,
			u.config.Kafka.IssueOrderTopics.IssueOrder,
			issueOrderMessage,
		)

		err = u.service.OrderHistory().UndoIssueOrder(ctx, u.repository.OrderHistory(), issueOrderMessage.Order.Id)
		if err != nil {
			u.metrics.Error()
			if errors.Is(err, models.RetryError) {
				err = u.RetryUndoIssueOrder(issueOrderMessage)
				if err != nil {
					u.metrics.KafkaError()
					log.Println(err)
				} else {
					log.Printf("consumer %s: -> %s: %v",
						u.config.Application.Name,
						u.config.Kafka.IssueOrderTopics.UndoIssueOrder,
						issueOrderMessage,
					)
				}
			} else {
				u.metrics.KafkaError()
				log.Println(err)
			}
			continue
		}

		log.Printf("consumer %s: <- %s: done",
			u.config.Application.Name,
			u.config.Kafka.IssueOrderTopics.UndoIssueOrder,
		)
	}
	return nil
}

func (u *UndoIssueOrderHandler) RetryUndoIssueOrder(message kafka.IssueOrderMessage) error {
	message.Base.SenderServiceName = u.config.Application.Name
	message.Base.Attempt += 1
	maxAttempts := int64(5)

	if message.Base.Attempt > maxAttempts {
		return models.MaxAttemptsError(nil, maxAttempts)
	}

	part, offs, kerr, err := kafka.SendMessage(u.producer, u.config.Kafka.IssueOrderTopics.UndoIssueOrder, message)
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
