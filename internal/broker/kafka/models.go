package kafka

type Base struct {
	SenderServiceName string
	Attempt           int64
	StartedAt         int64
}

type Order struct {
	Id int64
}

type Address struct {
	Id int64
}

type IssueOrderMessage struct {
	Base    Base
	Order   Order
	Address Address
}
