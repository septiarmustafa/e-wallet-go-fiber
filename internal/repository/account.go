package repository

import (
	"context"
	"database/sql"
	"e-wallet/domain"

	"github.com/doug-martin/goqu/v9"
)

type accountRespository struct {
	db *goqu.Database
}

func NewAccount(con *sql.DB) domain.AccountRepository {
	return &accountRespository{
		db: goqu.New("default", con),
	}
}

// FindByUserID implements domain.AccountRepository.
func (a *accountRespository) FindByUserID(ctx context.Context, id int64) (account domain.Account, err error) {
	dataset := a.db.From("accounts").Where(goqu.Ex{
		"user_id": id,
	})
	_, err = dataset.ScanStructContext(ctx, &account)
	return
}

// FindByAccountNumber implements domain.AccountRepository.
func (a *accountRespository) FindByAccountNumber(ctx context.Context, accNumber string) (account domain.Account, err error) {
	dataset := a.db.From("accounts").Where(goqu.Ex{
		"account_number": accNumber,
	})
	_, err = dataset.ScanStructContext(ctx, &account)
	return
}

// Update implements domain.AccountRepository.
func (a *accountRespository) Update(ctx context.Context, account *domain.Account) error {
	executor := a.db.Update("accounts").Where(goqu.Ex{
		"id": account.ID,
	}).Set(goqu.Record{
		"balance": account.Balance,
	}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
