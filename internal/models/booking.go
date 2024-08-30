package models

import "time"

type Booking struct {
	ID        int    `db:"id"`
	RoomID    string `db:"room_id"`
	Room      `db:"room"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}

type BookingDTO struct {
	ID        int       `json:"id"`
	RoomID    string    `json:"room_id"`
	Room      RoomDTO   `json:"room"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type BookingCreate struct {
	RoomID    string `json:"room_id" validate:"required"`
	StartTime string `json:"start_time" validate:"required,datetime=2006-01-02-15:04:05"` // YYYY-MM-DD-hh:mm:ss format
	EndTime   string `json:"end_time" validate:"required,datetime=2006-01-02-15:04:05"`   // YYYY-MM-DD-hh:mm:ss format
}
