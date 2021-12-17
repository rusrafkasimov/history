package dto

type AccountHistory struct {
	ID                   string `json:"id"`
	OperationID          string `json:"operation_id"`
	ClientID             string `json:"client_id"`
	AccountType          int32  `json:"account_type"`
	OperationCode        int32  `json:"operation_code"`
	OperationDesc        string `json:"operation_desc"`
	TransactionCode      int32  `json:"transaction_code"`
	TransactionDesc      string `json:"transaction_desc"`
	Comment              string `json:"comment"`
	State                int32  `json:"state"`
	StateDesc            string `json:"state_desc"`
	ScoreSum             int64  `json:"score_sum"`
	AvailableBalance     uint64 `json:"available_balance"`
	TransactionSum       int64  `json:"transaction_sum"`
	TransactionTimestamp string `json:"transaction_timestamp"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	ArchivedAt           string `json:"archived_at"`
	DeletedAt            string `json:"deleted_at"`
} // @Name AccountHistory
