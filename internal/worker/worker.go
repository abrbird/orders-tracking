package worker

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	cnfg "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/config"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/broker/kafka"
	rpstr "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/repository"
	srvc "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/service"
	"log"
)

type OrdersTrackingWorker struct {
	config                   *cnfg.Config
	producer                 sarama.SyncProducer
	issueOrderConsumer       *IssueOrderHandler
	undoIssueOrderConsumer   *UndoIssueOrderHandler
	confirmIssueOrderHandler *ConfirmIssueOrderHandler
}

func New(cfg *cnfg.Config, repository rpstr.Repository, service srvc.Service) (*OrdersTrackingWorker, error) {

	brokerConfig := kafka.NewConfig()
	producer, err := kafka.NewSyncProducer(cfg.Kafka.Brokers.String(), brokerConfig)
	if err != nil {
		return nil, err
	}

	worker := &OrdersTrackingWorker{
		config:   cfg,
		producer: producer,
		issueOrderConsumer: &IssueOrderHandler{
			producer:   producer,
			repository: repository,
			service:    service,
			config:     cfg,
		},
		undoIssueOrderConsumer: &UndoIssueOrderHandler{
			producer:   producer,
			repository: repository,
			service:    service,
			config:     cfg,
		},
		confirmIssueOrderHandler: &ConfirmIssueOrderHandler{
			producer:   producer,
			repository: repository,
			service:    service,
			config:     cfg,
		},
	}

	return worker, nil
}

func (w *OrdersTrackingWorker) StartConsuming(ctx context.Context) error {

	brokerConfig := kafka.NewConfig()
	issueOrder, err := sarama.NewConsumerGroup(
		w.config.Kafka.Brokers.String(),
		fmt.Sprintf("%s%sCG", w.config.Application.Name, w.config.Kafka.IssueOrderTopics.IssueOrder),
		brokerConfig,
	)
	if err != nil {
		return err
	}

	undoIssueOrder, err := sarama.NewConsumerGroup(
		w.config.Kafka.Brokers.String(),
		fmt.Sprintf("%s%sCG", w.config.Application.Name, w.config.Kafka.IssueOrderTopics.UndoIssueOrder),
		brokerConfig,
	)
	if err != nil {
		return err
	}

	confirmIssueOrder, err := sarama.NewConsumerGroup(
		w.config.Kafka.Brokers.String(),
		fmt.Sprintf("%s%sCG", w.config.Application.Name, w.config.Kafka.IssueOrderTopics.ConfirmIssueOrder),
		brokerConfig,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			err := issueOrder.Consume(ctx, []string{w.config.Kafka.IssueOrderTopics.IssueOrder}, w.issueOrderConsumer)
			if err != nil {
				log.Printf("%s consumer error: %v", w.config.Kafka.IssueOrderTopics.IssueOrder, err)
				//time.Sleep(time.Second * 5)
			}
		}
	}()
	go func() {
		for err := range issueOrder.Errors() {
			log.Println(err)
		}
	}()

	go func() {
		for {
			err := undoIssueOrder.Consume(ctx, []string{w.config.Kafka.IssueOrderTopics.UndoIssueOrder}, w.undoIssueOrderConsumer)
			if err != nil {
				log.Printf("%s consumer error: %v", w.config.Kafka.IssueOrderTopics.UndoIssueOrder, err)
				//time.Sleep(time.Second * 5)
			}
		}
	}()
	go func() {
		for err := range undoIssueOrder.Errors() {
			log.Println(err)
		}
	}()

	go func() {
		for {
			err := confirmIssueOrder.Consume(ctx, []string{w.config.Kafka.IssueOrderTopics.ConfirmIssueOrder}, w.confirmIssueOrderHandler)
			if err != nil {
				log.Printf("%s consumer error: %v", w.config.Kafka.IssueOrderTopics.ConfirmIssueOrder, err)
				//time.Sleep(time.Second * 5)
			}
		}
	}()
	go func() {
		for err := range confirmIssueOrder.Errors() {
			log.Println(err)
		}
	}()

	return nil
}
