package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dipper_backend/internal/model"
)

type Users interface {
	CreateUser(user model.User) error
	GetUserByTelegramID(telegramID string) (model.User, error)
	UpdateUserBalance(coins int, action string, telegramID string) error
	UpdateUserEnergy(energy int, action string, telegramID string) error
	ImproveUser(addedLevel int, addedMaxEnergy int, addedCoinsPerTap int, addedPassiveIncome int, telegramID string) error
	GetRefUsers(telegramID string) ([]model.User, error)
	SetLastEnergyUpdate(telegramID string, lastUpdate time.Time) error
}

type Upgrades interface {
	BuyUpgrade(userID int, upgradeID int) error
	GetAllUpgrades() ([]model.Upgrade, error)
	GetUserUpgrades(userID int) ([]model.Upgrade, error)
	CreateNewUpgrade(upgrade model.Upgrade) (int, error)
	GetOneUpgrade(upgradeID int) (model.Upgrade, error)
}

type Notifications interface {
	CreateNotification(notification model.Notification) error
	GetUserNotifications(userID int) ([]model.Notification, error)
}

type Repository struct {
	Users
	Upgrades
	Notifications
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Users: NewUsersPostgres(db), Upgrades: NewUpgradesPostgres(db), Notifications: NewNotificationsPostgres(db)}
}
