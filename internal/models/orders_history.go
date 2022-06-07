package models

const (
	Created       string = "created"
	Paid                 = "paid"
	Shipping             = "shipping"
	OnIssuePoint         = "on_issue_point"
	ReadyForIssue        = "ready_for_issue"
	Issued               = "issued"
)

const (
	InProgress string = "in_progress"
	Declined          = "declined"
	Confirmed         = "confirmed"
)

type OrderHistoryRecord struct {
	Id           int64
	OrderId      int64
	Status       string
	Confirmation string
}

type OrderHistoryRecordRetrieve struct {
	OrderHistoryRecord *OrderHistoryRecord
	Error              error
}

type OrderHistory []OrderHistoryRecord

type OrderHistoryRetrieve struct {
	OrderHistory *OrderHistory
	Error        error
}
