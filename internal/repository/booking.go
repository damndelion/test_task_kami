package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/damndelion/test_task_kami/internal/customrErrors"
	"github.com/damndelion/test_task_kami/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type IBookingRepo interface {
	CreateReservation(ctx context.Context, input models.Booking) (int, error)
	GetReservationsByRoomID(ctx context.Context, roomID string) ([]models.Booking, error)
	CheckOverlappingReservation(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error)
}

type BookingRepo struct {
	db *pgxpool.Pool
}

func NewBookingRepo(db *pgxpool.Pool) IBookingRepo {
	return &BookingRepo{db: db}
}

func (r *BookingRepo) CreateReservation(ctx context.Context, input models.Booking) (int, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Double check
	buildQuery := queryBuilder.Select("COUNT(*)").
		From("bookings").
		Where(squirrel.Eq{"room_id": input.RoomID}).
		Where("start_time < ? AND end_time > ?", input.EndTime, input.StartTime)

	query, args, err := buildQuery.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build check query: %w", err)
	}

	var count int
	err = tx.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to check overlapping reservations in transaction: %w", err)
	}
	if count > 0 {
		return 0, customrErrors.ErrAlreadyBooked
	}

	// Insert the reservation
	insertQueryBuilder := queryBuilder.Insert("bookings").
		Columns("room_id", "start_time", "end_time").
		Values(input.RoomID, input.StartTime, input.EndTime).
		Suffix("RETURNING id")

	insertSql, insertArgs, err := insertQueryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)
	}

	var id int
	err = tx.QueryRow(ctx, insertSql, insertArgs...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create reservation: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
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

func (r *BookingRepo) CheckOverlappingReservation(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error) {
	buildQuery := queryBuilder.Select("COUNT(*)").
		From("bookings").
		Where(squirrel.Eq{"room_id": roomID}).
		Where("start_time < ? AND end_time > ?", endTime, startTime)

	query, args, err := buildQuery.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build SQL query: %w", err)
	}

	var count int
	err = r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check overlapping reservations: %w", err)
	}

	return count > 0, nil
}
