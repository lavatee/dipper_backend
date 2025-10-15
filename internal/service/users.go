package service

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/lavatee/dipper_backend/internal/model"
	"github.com/lavatee/dipper_backend/internal/repository"
)

type Level struct {
	Goal             int
	AddedCoinsPerTap int
	AddedMaxEnergy   int
}

var levels = map[int]Level{
	1: {Goal: 100, AddedCoinsPerTap: 1},
	2: {Goal: 1000, AddedCoinsPerTap: 2},
	3: {Goal: 5000, AddedCoinsPerTap: 1},
	4: {Goal: 25000, AddedCoinsPerTap: 3},
	5: {Goal: 50000, AddedCoinsPerTap: 2},
	6: {Goal: 100000},
}

const (
	firstUserEnergy    = 100
	firstMaxUserEnergy = 100
	batchSize          = 5
	maxLevel           = 6
)

type UsersService struct {
	repo *repository.Repository
}

func NewUsersService(repo *repository.Repository) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s *UsersService) Login(user model.User) (model.User, error) {
	userInfo, err := s.repo.Users.GetUserByTelegramID(user.TelegramID)
	if err == nil {
		newEnergy, err := s.UpdateUserEnergy(userInfo)
		if err != nil {
			return model.User{}, err
		}
		userInfo.Energy += newEnergy
		return userInfo, nil
	}
	ref, err := uuid.GenerateUUID()
	if err != nil {
		return model.User{}, err
	}
	newUser := model.User{
		TelegramID:       user.TelegramID,
		TelegramUsername: user.TelegramUsername,
		UserRef:          ref,
		ByRef:            user.ByRef,
		Energy:           firstUserEnergy,
		MaxEnergy:        firstMaxUserEnergy,
	}
	if err := s.repo.Users.CreateUser(newUser); err != nil {
		return model.User{}, err
	}
	user, err = s.repo.Users.GetUserByTelegramID(user.TelegramID)
	return user, err
}

func (s *UsersService) UpdateUserEnergy(user model.User) (int, error) {
	thisUpdateTime := time.Now()
	lastUpdate := user.LastEnergyUpdate
	if lastUpdate.IsZero() {
		lastUpdate = thisUpdateTime
	}
	seconds := int(time.Since(lastUpdate).Seconds())
	restoredEnergy := int(seconds / 10)
	if restoredEnergy <= 0 {
		return 0, nil
	}
	newEnergy := user.Energy + restoredEnergy
	if newEnergy > user.MaxEnergy {
		restoredEnergy = user.MaxEnergy - user.Energy
		if restoredEnergy <= 0 {
			return 0, nil
		}
	}
	if err := s.repo.Users.UpdateUserEnergy(restoredEnergy, "+", user.TelegramID); err != nil {
		return 0, err
	}
	return restoredEnergy, s.repo.Users.SetLastEnergyUpdate(user.TelegramID, thisUpdateTime)
}

func (s *UsersService) ImproveUserByCoins(user model.User, coinsAmount int) error {
	addedLevel := 0
	addedMaxEnergy := 0
	addedCoinsPerTap := 0
	if user.Level < maxLevel && user.CoinsBalance+coinsAmount >= levels[user.Level].Goal {
		addedLevel = 1
		addedMaxEnergy = 200
		addedCoinsPerTap = levels[user.Level].AddedCoinsPerTap
	}
	return s.repo.Users.ImproveUser(addedLevel, addedMaxEnergy, addedCoinsPerTap, 0, user.TelegramID)
}

func (s *UsersService) TapsBatch(telegramID string) error {
	user, err := s.repo.Users.GetUserByTelegramID(telegramID)
	if err != nil {
		return err
	}
	coinsAmount := user.CoinsPerTap * batchSize
	if user.Energy < batchSize {
		return fmt.Errorf("energy")
	}
	if err := s.ImproveUserByCoins(user, coinsAmount); err != nil {
		return err
	}
	if err := s.repo.Users.UpdateUserBalance(coinsAmount, "+", telegramID); err != nil {
		return err
	}
	if err := s.repo.Users.UpdateUserEnergy(batchSize, "-", telegramID); err != nil {
		return err
	}
	return nil
}

func (s *UsersService) GetRefUsers(telegramID string) ([]model.User, error) {
	return s.repo.Users.GetRefUsers(telegramID)
}
