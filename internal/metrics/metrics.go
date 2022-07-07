package metrics

type Metrics interface {
	Error()
	KafkaError()
}
