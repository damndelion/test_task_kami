package service

import (
	"context"
	"fmt"
	"github.com/damndelion/test_task_kami/internal/customrErrors"
	"github.com/damndelion/test_task_kami/internal/models"
	"github.com/damndelion/test_task_kami/internal/repository"
	"go.uber.org/zap"
	"time"
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
	//prepare data
	startTime, err := time.Parse("2006-01-02-15:04:05", input.StartTime)
	if err != nil {
		return 0, fmt.Errorf("invalid start time format: %w", err)
	}
	endTime, err := time.Parse("2006-01-02-15:04:05", input.EndTime)
	if err != nil {
		return 0, fmt.Errorf("invalid end time format: %w", err)
	}
	if startTime.After(endTime) {
		return 0, customrErrors.ErrEndBeforeStart
	}

	booking := models.Booking{
		RoomID:    input.RoomID,
		StartTime: startTime,
		EndTime:   endTime,
	}

	//check for overlapping
	overlapping, err := s.repo.CheckOverlappingReservation(ctx, booking.RoomID, booking.StartTime, booking.EndTime)
	if err != nil {
		s.logs.Named("booking").Errorf("booking - service - CheckOverlappingReservation - failed: %v", err)
		return 0, fmt.Errorf("booking - service - CheckOverlappingReservation - failed: %w", err)
	}
	if overlapping {
		return 0, customrErrors.ErrAlreadyBooked
	}

	// create reservation
	id, err := s.repo.CreateReservation(ctx, booking)
	if err != nil {
		s.logs.Named("booking").Errorf("booking - repo - CreateReservation - failed to create reservation: %v", err)
		return 0, fmt.Errorf("booking - repo - CreateReservation - failed to create reservation: %w", err)
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
