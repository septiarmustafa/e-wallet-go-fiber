package domain

import (
	"context"
	"e-wallet/dto"
	"time"
)

type Transaction struct {
	ID                  int64     `db:"id"`
	Account_id          int64     `db:"account_id"`
	SofNumber           string    `db:"sof_number"`
	DofNumber           string    `db:"dof_number"`
	TransactionType     string    `db:"transaction_type"`
	Amount              float64   `db:"amount"`
	TransactionDatetime time.Time `db:"transaction_datetime"`
}

type TransactionRepository interface {
	Insert(ctx context.Context, transaction *Transaction) error
}

type TransactionService interface {
	TransferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error)
	TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error
}
