package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/damndelion/test_task_kami/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IBookingRepo interface {
	CreateReservation(ctx context.Context, input models.BookingCreate) (int, error)
	GetReservationsByRoomID(ctx context.Context, roomID string) ([]models.Booking, error)
}

type BookingRepo struct {
	db *pgxpool.Pool
}

func NewBookingRepo(db *pgxpool.Pool) IBookingRepo {
	return &BookingRepo{db: db}
}

func (r *BookingRepo) CreateReservation(ctx context.Context, input models.BookingCreate) (int, error) {
	// build a sql query
	query := queryBuilder.
		Insert("bookings").
		Columns("room_id", "start_time", "end_time").
		Values(input.RoomID, input.StartTime, input.EndTime).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var id int
	// run the query
	err = r.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *BookingRepo) GetReservationsByRoomID(ctx context.Context, roomID string) ([]models.Booking, error) {
	// Build a SQL query
	query := queryBuilder.
		Select(`b.id, b.room_id, b.start_time, b.end_time,
			r.room_name, r.id`).
		From("bookings as b").
		Join("rooms as r ON r.id = b.room_id").
		Where(squirrel.Eq{"b.room_id": roomID}).
		OrderBy("b.start_time DESC")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	// Run the query
	rows, err := r.db.Query(ctx, rawQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Booking
	for rows.Next() {
		var booking models.Booking
		var roomName string
		var roomIDScan string

		err = rows.Scan(
			&booking.ID,
			&booking.RoomID,
			&booking.StartTime,
			&booking.EndTime,
			&roomName,
			&roomIDScan,
		)
		if err != nil {
			return nil, err
		}

		booking.Room.RoomName = roomName
		booking.Room.ID = roomIDScan

		res = append(res, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
