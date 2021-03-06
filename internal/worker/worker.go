package worker

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	cnfg "github.com/abrbird/orders-tracking/config"
	"github.com/abrbird/orders-tracking/internal/broker/kafka"
	"github.com/abrbird/orders-tracking/internal/metrics"
	rpstr "github.com/abrbird/orders-tracking/internal/repository"
	srvc "github.com/abrbird/orders-tracking/internal/service"
	"log"
)

type OrdersTrackingWorker struct {
	config                   *cnfg.Config
	repository               rpstr.Repository
	service                  srvc.Service
	metrics                  metrics.Metrics
	producer                 sarama.SyncProducer
	issueOrderConsumer       *IssueOrderHandler
	undoIssueOrderConsumer   *UndoIssueOrderHandler
	confirmIssueOrderHandler *ConfirmIssueOrderHandler
}

func New(cfg *cnfg.Config, repository rpstr.Repository, service srvc.Service, metrics metrics.Metrics) (*OrdersTrackingWorker, error) {

	brokerConfig := kafka.NewConfig()
	producer, err := kafka.NewSyncProducer(cfg.Kafka.Brokers.String(), brokerConfig)
	if err != nil {
		return nil, err
	}

	worker := &OrdersTrackingWorker{
		config:     cfg,
		repository: repository,
		service:    service,
		metrics:    metrics,
		producer:   producer,
		issueOrderConsumer: &IssueOrderHandler{
			producer:   producer,
			repository: repository,
			service:    service,
			metrics:    metrics,
			config:     cfg,
		},
		undoIssueOrderConsumer: &UndoIssueOrderHandler{
			producer:   producer,
			repository: repository,
			service:    service,
			metrics:    metrics,
			config:     cfg,
		},
		confirmIssueOrderHandler: &ConfirmIssueOrderHandler{
			producer:   producer,
			repository: repository,
			service:    service,
			metrics:    metrics,
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
