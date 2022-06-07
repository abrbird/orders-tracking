package worker

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	cnfg "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/config"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/broker/kafka"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
	rpstr "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/repository"
	srvc "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/service"
	"log"
)

type ConfirmIssueOrderHandler struct {
	producer   sarama.SyncProducer
	repository rpstr.Repository
	service    srvc.Service
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

		log.Printf("consumer %s: -> %s: %v",
			i.config.Application.Name,
			i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder,
			msg.Value,
		)

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
			log.Print("Unmarshall failed: value=%v, err=%v", string(msg.Value), err)
			continue
		}

		ctx := context.Background()
		lastOrderHistoryRecordRetrieved := i.service.OrderHistory().RetrieveByStatus(
			ctx,
			i.repository.OrderHistory(),
			issueOrderMessage.Order.Id,
			models.ReadyForIssue,
		)

		if lastOrderHistoryRecordRetrieved.Error != nil {
			log.Printf("error on message processing: %v", err)
			i.RetryConfirmIssueOrder(issueOrderMessage)
			continue
		}

		if lastOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation == models.Confirmed {
			log.Printf("order is already issued: %v", err)
			continue
		}

		confirmIssueRecord := models.OrderHistoryRecord{
			Id:           lastOrderHistoryRecordRetrieved.OrderHistoryRecord.Id,
			OrderId:      lastOrderHistoryRecordRetrieved.OrderHistoryRecord.OrderId,
			Status:       lastOrderHistoryRecordRetrieved.OrderHistoryRecord.Status,
			Confirmation: models.InProgress,
		}

		err = i.service.OrderHistory().Update(
			ctx,
			i.repository.OrderHistory(),
			&confirmIssueRecord,
		)
		if err != nil {
			log.Printf("order can not be issued: %v", err)
			i.RetryConfirmIssueOrder(issueOrderMessage)
			continue
		}

		log.Printf("consumer %s: <- %s: done",
			i.config.Application.Name,
			i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder,
		)
	}
	return nil
}

func (i *ConfirmIssueOrderHandler) RetryConfirmIssueOrder(message kafka.IssueOrderMessage) {
	message.Base.SenderServiceName = i.config.Application.Name
	message.Base.Attempt += 1

	part, offs, kerr, err := kafka.SendMessage(i.producer, i.config.Kafka.IssueOrderTopics.IssueOrder, message)
	if err != nil {
		log.Printf("can not send message: %v", err)
		return
	}

	if kerr != nil {
		log.Printf("can not send message: %v", kerr)
		return
	}

	log.Printf("consumer %s: %v -> %v", i.config.Application.Name, part, offs)
	return
}
