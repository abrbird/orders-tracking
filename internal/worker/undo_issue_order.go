package worker

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	cnfg "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/config"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/broker/kafka"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/models"
	"log"

	rpstr "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/repository"
	srvc "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/service"
)

type UndoIssueOrderHandler struct {
	producer   sarama.SyncProducer
	repository rpstr.Repository
	service    srvc.Service
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
			log.Print("Unmarshall failed: value=%v, err=%v", string(msg.Value), err)
			continue
		}

		log.Printf("consumer %s: <- %s: %v",
			u.config.Application.Name,
			u.config.Kafka.IssueOrderTopics.IssueOrder,
			issueOrderMessage,
		)

		ctx := context.Background()
		issuedOrderHistoryRecordRetrieved := u.service.OrderHistory().RetrieveByStatus(
			ctx,
			u.repository.OrderHistory(),
			issueOrderMessage.Order.Id,
			models.Issued,
		)

		if issuedOrderHistoryRecordRetrieved.Error != nil {
			log.Printf("error on message processing: %v", err)
			u.RetryUndoIssueOrder(issueOrderMessage)
			continue
		}

		if issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation == models.InProgress {
			issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation = models.Declined

			err = u.service.OrderHistory().Update(
				ctx,
				u.repository.OrderHistory(),
				issuedOrderHistoryRecordRetrieved.OrderHistoryRecord,
			)
			if err != nil {
				log.Printf("order can not be issued: %v", err)
				u.RetryUndoIssueOrder(issueOrderMessage)
				continue
			}
		}
		log.Printf("consumer %s: <- %s: done",
			u.config.Application.Name,
			u.config.Kafka.IssueOrderTopics.UndoIssueOrder,
		)
	}
	return nil
}

func (u *UndoIssueOrderHandler) RetryUndoIssueOrder(message kafka.IssueOrderMessage) {
	message.Base.SenderServiceName = u.config.Application.Name
	message.Base.Attempt += 1

	if message.Base.Attempt > 5 {
		log.Printf("reached max attempts: %v", message)
		return
	}

	part, offs, kerr, err := kafka.SendMessage(u.producer, u.config.Kafka.IssueOrderTopics.UndoIssueOrder, message)
	if err != nil {
		log.Printf("can not send message: %v", err)
		return
	}

	if kerr != nil {
		log.Printf("can not send message: %v", kerr)
		return
	}

	log.Printf("consumer %s: sent %v -> %v", u.config.Application.Name, part, offs)
	return
}
