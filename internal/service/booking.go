package service

import (
	"context"
	"fmt"
	"github.com/damndelion/test_task_kami/internal/models"
	"github.com/damndelion/test_task_kami/internal/repository"
	"go.uber.org/zap"
)

type IBookingService interface {
	CreateReservation(ctx context.Context, input models.BookingCreate) (int, error)
	GetReservationByRoomID(ctx context.Context, roomID string) ([]models.BookingDTO, error)
}

type BookingService struct {
	logs *zap.SugaredLogger
	repo repository.IBookingRepo
}

func NewBookingService(logs *zap.SugaredLogger, repo repository.IBookingRepo) IBookingService {
	return &BookingService{
		logs: logs,
		repo: repo,
	}
}

func (s *BookingService) CreateReservation(ctx context.Context, input models.BookingCreate) (int, error) {
	id, err := s.repo.CreateReservation(ctx, input)
	if err != nil {
		s.logs.Named("booking").Errorf("booking  - repo - CreateReservation - failed to create reservation: %v", err)
		return 0, fmt.Errorf("booking  - repo - CreateReservation - failed to create reservation: %w", err)
	}

	return id, nil
}

func (s *BookingService) GetReservationByRoomID(ctx context.Context, roomID string) ([]models.BookingDTO, error) {
	res, err := s.repo.GetReservationsByRoomID(ctx, roomID)
	if err != nil {
		s.logs.Named("booking").Errorf("booking  - repo - GetReservationsByRoomID - failed to get reservations: %v", err)
		return nil, fmt.Errorf("booking  - repo - GetReservationsByRoomID - failed to get reservations: %w", err)
	}
	resDTOs := MapBookingsToDTOs(res)

	return resDTOs, nil
}

func MapBookingsToDTOs(bookings []models.Booking) []models.BookingDTO {
	dto := make([]models.BookingDTO, len(bookings))
	for i, booking := range bookings {
		dto[i] = MapBookingToDTO(booking)
	}
	return dto
}

func MapBookingToDTO(booking models.Booking) models.BookingDTO {
	return models.BookingDTO{
		ID:     booking.ID,
		RoomID: booking.RoomID,
		Room: models.RoomDTO{
			ID:       booking.Room.ID,
			RoomName: booking.Room.RoomName,
		},
		StartTime: booking.StartTime,
		EndTime:   booking.EndTime,
	}
}
