package service

import (
	"github.com/damndelion/test_task_kami/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	IBookingService
}

func NewService(logs *zap.SugaredLogger, repo *repository.Repository) *Service {
	return &Service{
		IBookingService: NewBookingService(logs, repo.IBookingRepo),
	}
}
