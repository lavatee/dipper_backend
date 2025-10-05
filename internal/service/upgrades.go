package service

import (
	"fmt"

	"github.com/lavatee/dipper_backend/internal/model"
	"github.com/lavatee/dipper_backend/internal/repository"
)

type UpgradesService struct {
	repo *repository.Repository
}

func NewUpgradesService(repo *repository.Repository) *UpgradesService {
	return &UpgradesService{
		repo: repo,
	}
}

var admins = map[string]bool{
	"123": true,
}

func (s *UpgradesService) CreateUpgrade(upgrade model.Upgrade, creatorTelegramID string) (int, error) {
	if _, exists := admins[creatorTelegramID]; !exists {
		return 0, fmt.Errorf("notadmin")
	}
	return s.repo.Upgrades.CreateNewUpgrade(upgrade)
}

func (s *UpgradesService) GetAllUpgrades() ([]model.Upgrade, error) {
	return s.repo.Upgrades.GetAllUpgrades()
}

func (s *UpgradesService) GetUserUpgrades(telegramID string) ([]model.Upgrade, error) {
	user, err := s.repo.Users.GetUserByTelegramID(telegramID)
	if err != nil {
		return nil, err
	}
	return s.repo.Upgrades.GetUserUpgrades(user.ID)
}

func (s *UpgradesService) BuyUpgradeByCoins(telegramID string, upgradeID int) error {
	upgrade, err := s.repo.Upgrades.GetOneUpgrade(upgradeID)
	if err != nil {
		return err
	}
	user, err := s.repo.Users.GetUserByTelegramID(telegramID)
	if err != nil {
		return err
	}
	if err := s.repo.Users.UpdateUserBalance(upgrade.CoinsPrice, "-", telegramID); err != nil {
		return err
	}
	if err := s.repo.Upgrades.BuyUpgrade(user.ID, upgradeID); err != nil {
		return err
	}
	if err := s.repo.Users.ImproveUser(0, upgrade.NewMaxEnergy, upgrade.NewOneTapCoins, upgrade.NewPassiveIncome, telegramID); err != nil {
		return err
	}
	return nil
}
