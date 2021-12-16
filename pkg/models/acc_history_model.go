package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type AccountHistory struct {
	ID                   uuid.UUID `json:"id"`
	OperationID          string    `json:"operation_id"`
	ClientID             string    `json:"client_id"`
	AccountType          int32     `json:"account_type"`
	OperationCode        int32     `json:"operation_code"`
	OperationDesc        string    `json:"operation_desc"`
	TransactionCode      int32     `json:"transaction_code"`
	TransactionDesc      string    `json:"transaction_desc"`
	Comment              string    `gorm:"type:text" json:"comment"`
	State                int32     `json:"state"`
	StateDesc            string    `json:"state_desc"`
	ScoreSum             int64     `json:"score_sum"`
	AvailableBalance     uint64    `json:"available_balance"`
	TransactionSum       int64     `json:"transaction_sum"`
	TransactionTimestamp time.Time `json:"transaction_timestamp"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	ArchivedAt           time.Time `json:"archived_at"`
	DeletedAt            time.Time `json:"deleted_at"`
}
