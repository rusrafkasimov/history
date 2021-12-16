package convert

import (
	"github.com/rusrafkasimov/history/pkg/dto"
	"github.com/rusrafkasimov/history/pkg/models"
)

func AccountHistoryDto (history dto.AccountHistory) (models.AccountHistory, error) {
	return models.AccountHistory{
		OperationID: history.OperationID,
		ClientID: history.ClientID,
		AccountType: history.AccountType,
		OperationCode: history.OperationCode,
		OperationDesc: history.OperationDesc,
		TransactionCode: history.TransactionCode,
		TransactionDesc: history.TransactionDesc,
		Comment: history.Comment,
		State: history.State,
		StateDesc: history.StateDesc,
		ScoreSum: history.ScoreSum,
		AvailableBalance: history.AvailableBalance,
		TransactionSum: history.TransactionSum,
	}, nil
}

func AccountHistoryModel(history models.AccountHistory) dto.AccountHistory {
	return dto.AccountHistory{
		ID: history.ID.String(),
		OperationID: history.OperationID,
		ClientID: history.ClientID,
		AccountType: history.AccountType,
		OperationCode: history.OperationCode,
		OperationDesc: history.OperationDesc,
		TransactionCode: history.TransactionCode,
		TransactionDesc: history.TransactionDesc,
		Comment: history.Comment,
		State: history.State,
		StateDesc: history.StateDesc,
		ScoreSum: history.ScoreSum,
		AvailableBalance: history.AvailableBalance,
		TransactionSum: history.TransactionSum,
		TransactionTimestamp: history.TransactionTimestamp.String(),
		CreatedAt: history.CreatedAt.String(),
		UpdatedAt: history.UpdatedAt.String(),
		ArchivedAt: history.ArchivedAt.String(),
		DeletedAt: history.DeletedAt.String(),
	}
}