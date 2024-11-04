package repository

import (
	"context"
	"database/sql"
	"e-wallet/domain"

	"github.com/doug-martin/goqu/v9"
)

type UserRepository struct {
	db *goqu.Database
}

func NewUser(con *sql.DB) domain.UserRepository {
	return &UserRepository{
		db: goqu.New("default", con),
	}
}

// FindByID implements domain.UserRepository.
func (u *UserRepository) FindByID(ctx context.Context, id int64) (user domain.User, err error) {
	dataset := u.db.From("users").Where(goqu.Ex{
		"id": id,
	})
	_, err = dataset.ScanStructContext(ctx, &user)
	return
}

// FindByUsername implements domain.UserRepository.
func (u *UserRepository) FindByUsername(ctx context.Context, username string) (user domain.User, err error) {
	dataset := u.db.From("users").Where(goqu.Ex{
		"username": username,
	})
	_, err = dataset.ScanStructContext(ctx, &user)
	return
}

// Insert implements domain.UserRepository.
func (u *UserRepository) Insert(ctx context.Context, user *domain.User) error {
	executor := u.db.Insert("users").Rows(goqu.Record{
		"full_name": user.FullName,
		"email":     user.Email,
		"phone":     user.Phone,
		"username":  user.Username,
		"password":  user.Password,
	}).Returning("id").Executor()
	_, err := executor.ScanStructContext(ctx, user)
	return err
}

// Update implements domain.UserRepository.
func (u *UserRepository) Update(ctx context.Context, user *domain.User) error {
	user.EmailVerifiedAtDB = sql.NullTime{
		Time:  user.EmailVerifiedAt,
		Valid: true,
	}

	executor := u.db.Update("users").Where(goqu.Ex{
		"id": user.ID,
	}).Set(goqu.Record{
		"full_name":         user.FullName,
		"email":             user.Email,
		"phone":             user.Phone,
		"username":          user.Username,
		"password":          user.Password,
		"email_verified_at": user.EmailVerifiedAtDB,
	}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
