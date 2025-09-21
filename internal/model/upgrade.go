package model

type Upgrade struct {
	ID               int    `json:"id" db:"id"`
	Name             string `json:"name" db:"name"`
	Description      string `json:"description" db:"description"`
	CoinsPrice       int    `json:"coins_price" db:"coins_price"`
	StarsPrice       int    `json:"stars_price" db:"stars_price"`
	NewMaxEnergy     int    `json:"new_max_energy" db:"new_max_energy"`
	NewOneTapCoins   int    `json:"new_one_tap_coins" db:"new_one_tap_coins"`
	NewPassiveIncome int    `json:"new_passive_income" db:"new_passive_income"`
}

type UserUpgrade struct {
	ID        int `json:"id" db:"id"`
	UserID    int `json:"user_id" db:"user_id"`
	UpgradeID int `json:"upgrade_id" db:"upgrade_id"`
}
