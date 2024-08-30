package repository

import "github.com/jackc/pgx/v5/pgxpool"

type IRoomRepo interface {
}

type RoomRepo struct {
	db *pgxpool.Pool
}

func NewRoomRepo(db *pgxpool.Pool) IRoomRepo {
	return &RoomRepo{db: db}
}
