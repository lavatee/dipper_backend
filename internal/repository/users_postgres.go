package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dipper_backend/internal/model"
)

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{db: db}
}

func (r *UsersPostgres) CreateUser(user model.User) error {
	query := fmt.Sprintf("INSERT INTO %s (telegram_username, telegram_id, energy, max_energy, user_ref, by_ref) VALUES ($1, $2, $3, $4, $5, $6)", usersTable)
	_, err := r.db.Exec(query, user.TelegramUsername, user.TelegramID, user.Energy, user.MaxEnergy, user.UserRef, user.ByRef)
	return err
}

func (r *UsersPostgres) GetUserByTelegramID(telegramID string) (model.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE telegram_id = $1", usersTable)
	var user model.User
	if err := r.db.Get(&user, query, telegramID); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UsersPostgres) UpdateUserBalance(coins int, action string, telegramID string) error { //action: "+" or "-"
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s SET coins_balance = coins_balance %s $1 WHERE telegram_id = $2 RETURNING coins_balance", usersTable, action)
	row := tx.QueryRow(query, coins, telegramID)
	var balance int
	if err := row.Scan(&balance); err != nil {
		tx.Rollback()
		return err
	}
	if balance < 0 {
		tx.Rollback()
		return fmt.Errorf("balance")
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *UsersPostgres) UpdateUserEnergy(energy int, action string, telegramID string) error {
	query := fmt.Sprintf("UPDATE %s SET energy = energy %s $1 AND last_energy_update = CURRENT_TIMESTAMP WHERE telegram_id = $2", usersTable, action)
	_, err := r.db.Exec(query, energy, telegramID)
	return err
}

func (r *UsersPostgres) ImproveUser(addedLevel int, addedMaxEnergy int, addedCoinsPerTap int, addedPassiveIncome int, telegramID string) error {
	query := fmt.Sprintf("UPDATE %s SET level = level + $1, max_energy = max_energy + $2, coins_per_tap = coins_per_tap + $3, passive_income = passive_income + $4 WHERE telegram_id = $5", usersTable)
	_, err := r.db.Exec(query, addedLevel, addedMaxEnergy, addedCoinsPerTap, addedPassiveIncome, telegramID)
	return err
}

func (r *UsersPostgres) GetRefUsers(telegramID string) ([]model.User, error) {
	var users []model.User
	query := fmt.Sprintf("SELECT telegram_username, telegram_id FROM %s WHERE by_ref = (SELECT user_ref FROM %s WHERE telegram_id = $1)", usersTable, usersTable)
	if err := r.db.Select(&users, query, telegramID); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UsersPostgres) SetLastEnergyUpdate(telegramID string, lastUpdate time.Time) error {
	query := fmt.Sprintf("UPDATE %s SET last_energy_update = $1 WHERE telegram_id = $2", usersTable)
	_, err := r.db.Exec(query, lastUpdate, telegramID)
	return err
}
