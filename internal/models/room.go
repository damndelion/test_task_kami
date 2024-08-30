package models

type Room struct {
	ID       string `db:"id"`
	RoomName string `db:"room_name"`
}

type RoomDTO struct {
	ID       string `json:"id"`
	RoomName string `json:"room_name"`
}
