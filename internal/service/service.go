package service

import (
	"github.com/damndelion/test_task_kami/internal/repository"
)

type Service struct {
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
