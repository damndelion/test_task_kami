package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

var queryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Repository struct {
	IBookingRepo
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		IBookingRepo: NewBookingRepo(db),
	}
}
