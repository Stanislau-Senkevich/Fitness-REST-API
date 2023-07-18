package postgres

import "github.com/jmoiron/sqlx"

type PartnershipRepository struct {
	db *sqlx.DB
}

func NewPartnershipRepository(db *sqlx.DB) *PartnershipRepository {
	return &PartnershipRepository{db: db}
}
