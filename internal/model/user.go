package model

import "time"

type User struct {
	ID                      int       `json:"id" db:"id"`
	TelegramUsername        string    `json:"telegram_username" db:"telegram_username"`
	TelegramID              string    `json:"telegram_id" db:"telegram_id"`
	CoinsBalance            int       `json:"coins_balance" db:"coins_balance"`
	CoinsPerTap             int       `json:"coins_per_tap" db:"coins_per_tap"`
	Level                   int       `json:"level" db:"level"`
	Energy                  int       `json:"energy" db:"energy"`
	MaxEnergy               int       `json:"max_energy" db:"max_energy"`
	LastEnergyUpdate        time.Time `json:"last_energy_update" db:"last_energy_update"`
	UserRef                 string    `json:"user_ref" db:"user_ref"`
	ByRef                   string    `json:"by_ref" db:"by_ref"`
	PassiveIncome           int       `json:"passive_income" db:"passive_income"`
	LastPassiveIncomeUpdate time.Time `json:"last_passive_income_update" db:"last_passive_income_update"`
}
