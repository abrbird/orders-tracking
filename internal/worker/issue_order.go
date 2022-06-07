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

type IssueOrderHandler struct {
	producer   sarama.SyncProducer
	repository rpstr.Repository
	service    srvc.Service
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

		log.Printf("consumer %s: <- %s: %v",
			i.config.Application.Name,
			i.config.Kafka.IssueOrderTopics.IssueOrder,
			msg.Value,
		)

		if msg.Topic != i.config.Kafka.IssueOrderTopics.IssueOrder {
			log.Printf(
				"topic names does not match: expected - %s, got %s\n",
				i.config.Kafka.IssueOrderTopics.IssueOrder,
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

		issuedOrderHistoryRecordRetrieved := i.service.OrderHistory().RetrieveByStatus(
			ctx,
			i.repository.OrderHistory(),
			issueOrderMessage.Order.Id,
			models.Issued,
		)
		if issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation == models.InProgress {
			log.Printf("order is on issuing: %v", err)
			i.RetryIssueOrder(issueOrderMessage)
			continue
		}
		if issuedOrderHistoryRecordRetrieved.OrderHistoryRecord.Confirmation == models.Confirmed {
			log.Printf("order is already issued: %v", err)
			continue
		}

		readyOrderHistoryRecordRetrieved := i.service.OrderHistory().RetrieveByStatus(
			ctx,
			i.repository.OrderHistory(),
			issueOrderMessage.Order.Id,
			models.ReadyForIssue,
		)

		if readyOrderHistoryRecordRetrieved.Error != nil {
			log.Printf("order can not be issued: %v", err)
			i.RetryIssueOrder(issueOrderMessage)
			continue
		}

		issuedOrderRecord := models.OrderHistoryRecord{
			OrderId:      readyOrderHistoryRecordRetrieved.OrderHistoryRecord.OrderId,
			Status:       models.Issued,
			Confirmation: models.InProgress,
		}

		err = i.service.OrderHistory().Create(
			ctx,
			i.repository.OrderHistory(),
			&issuedOrderRecord,
		)
		if err != nil {
			log.Printf("order can not be issued: %v", err)
			i.RetryIssueOrder(issueOrderMessage)
			continue
		}
		i.SendRemoveOrder(issueOrderMessage)

		log.Printf("consumer %s: -> %s: %v",
			i.config.Application.Name,
			i.config.Kafka.IssueOrderTopics.RemoveOrder,
			issueOrderMessage,
		)

	}
	return nil
}

func (i *IssueOrderHandler) RetryIssueOrder(message kafka.IssueOrderMessage) {
	message.Base.SenderServiceName = i.config.Application.Name
	message.Base.Attempt += 1

	if message.Base.Attempt > 5 {
		log.Printf("reached max attempts: %v", message)
		return
	}

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

func (i *IssueOrderHandler) SendRemoveOrder(message kafka.IssueOrderMessage) {
	message.Base.SenderServiceName = i.config.Application.Name

	part, offs, kerr, err := kafka.SendMessage(i.producer, i.config.Kafka.IssueOrderTopics.RemoveOrder, message)
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
