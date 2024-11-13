package repository

import (
	"context"
	"database/sql"
	"e-wallet/domain"

	"github.com/doug-martin/goqu/v9"
)

type templateRepository struct {
	db *goqu.Database
}

func NewTemplate(con *sql.DB) domain.TemplateRepository {
	return &templateRepository{
		db: goqu.New("default", con),
	}
}

// FindByCode implements domain.TemplateRepository.
func (t *templateRepository) FindByCode(ctx context.Context, code string) (temp domain.Template, err error) {
	dataset := t.db.From("templates").Where(goqu.Ex{
		"code": code,
	})

	_, err = dataset.ScanStructContext(ctx, &temp)
	return
}
