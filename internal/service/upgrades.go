package service

import "github.com/lavatee/dipper_backend/internal/repository"

type UpgradesService struct {
	repo *repository.Repository
}

func NewUpgradesService(repo *repository.Repository) *UpgradesService {
	return &UpgradesService{
		repo: repo,
	}
}
