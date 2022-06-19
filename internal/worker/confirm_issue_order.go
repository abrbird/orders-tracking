package worker

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	cnfg "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/config"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/broker/kafka"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/metrics"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
	rpstr "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/repository"
	srvc "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/service"
	"log"
)

type ConfirmIssueOrderHandler struct {
	producer   sarama.SyncProducer
	repository rpstr.Repository
	service    srvc.Service
	metrics    metrics.Metrics
	config     *cnfg.Config
}

func (i *ConfirmIssueOrderHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (i *ConfirmIssueOrderHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (i *ConfirmIssueOrderHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := context.Background()

		if msg.Topic != i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder {
			log.Printf(
				"topic names does not match: expected - %s, got %s\n",
				i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder,
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
			i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder,
			issueOrderMessage,
		)

		err = i.service.OrderHistory().ConfirmIssueOrder(ctx, i.repository.OrderHistory(), issueOrderMessage.Order.Id)
		if err != nil {
			i.metrics.Error()
			if errors.Is(err, models.RetryError) {
				err = i.RetryConfirmIssueOrder(issueOrderMessage)
				if err != nil {
					i.metrics.KafkaError()
					log.Println(err)
				} else {
					log.Printf("consumer %s: -> %s: %v",
						i.config.Application.Name,
						i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder,
						issueOrderMessage,
					)
				}
			} else {
				i.metrics.KafkaError()
				log.Println(err)
			}
			continue
		}

		log.Printf("consumer %s: <- %s: done",
			i.config.Application.Name,
			i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder,
		)
	}
	return nil
}

func (i *ConfirmIssueOrderHandler) RetryConfirmIssueOrder(message kafka.IssueOrderMessage) error {
	message.Base.SenderServiceName = i.config.Application.Name
	message.Base.Attempt += 1

	part, offs, kerr, err := kafka.SendMessage(i.producer, i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder, message)
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
