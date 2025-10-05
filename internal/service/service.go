package service

import (
	"github.com/lavatee/dipper_backend/internal/model"
	"github.com/lavatee/dipper_backend/internal/repository"
)

type Users interface {
	Login(user model.User) (model.User, error)
	UpdateUserEnergy(user model.User) (int, error)
	ImproveUserByCoins(user model.User, coinsAmount int) error
	TapsBatch(telegramID string) error
	GetRefUsers(telegramID string) ([]model.User, error)
}

type Upgrades interface {
	CreateUpgrade(upgrade model.Upgrade, creatorTelegramID string) (int, error)
	GetAllUpgrades() ([]model.Upgrade, error)
	GetUserUpgrades(telegramID string) ([]model.Upgrade, error)
	BuyUpgradeByCoins(telegramID string, upgradeID int) error
}

type Service struct {
	Users
	Upgrades
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Users:    NewUsersService(repo),
		Upgrades: NewUpgradesService(repo),
	}
}
