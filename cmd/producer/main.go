package main

import (
	"github.com/abrbird/orders-tracking/config"
	"github.com/abrbird/orders-tracking/internal/broker/kafka"
	"log"
	"time"
)

func main() {
	cfg, err := config.ParseConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	brokerCfg := kafka.NewConfig()
	syncProducer, err := kafka.NewSyncProducer(cfg.Kafka.Brokers.String(), brokerCfg)

	if err != nil {
		log.Fatal(err)
	}

	incomeServiceName := "Income"

	for {
		message := kafka.IssueOrderMessage{
			Base: kafka.Base{
				SenderServiceName: incomeServiceName,
				Attempt:           1,
				StartedAt:         time.Now().UnixNano(),
			},
			Order: kafka.Order{
				Id: 2,
			},
			Address: kafka.Address{
				Id: 1,
			},
		}

		part, offs, kerr, err := kafka.SendMessage(syncProducer, cfg.Kafka.IssueOrderTopics.IssueOrder, message)
		if err != nil {
			log.Printf("can not send message: %v", err)
			continue
		}

		log.Printf("producer %s: sent %v -> %v; err: %v", incomeServiceName, part, offs, kerr)
		break
		//time.Sleep(time.Millisecond * 5000)

		//if rand.Intn(10) == 9 {
		//	par, off, err = syncProducer.SendMessage(&sarama.ProducerMessage{
		//		Topic: cfg.Kafka.IssueOrderTopics.UndoIssueOrder,
		//		Key:   sarama.StringEncoder(fmt.Sprintf("%v", d.Id)),
		//		Value: sarama.ByteEncoder(b),
		//	})
		//	log.Printf("producer UndoIssueOrder %v -> %v; %v", par, off, err)
		//}
	}
}
