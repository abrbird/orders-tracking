package kafka

import "github.com/Shopify/sarama"

func NewConfig() *sarama.Config {
	brokerCfg := sarama.NewConfig()
	brokerCfg.Producer.Return.Successes = true
	//brokerCfg.Consumer.Offsets.AutoCommit.Enable = false
	return brokerCfg
}
