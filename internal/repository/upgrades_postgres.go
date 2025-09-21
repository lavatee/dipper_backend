package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dipper_backend/internal/model"
)

type UpgradesPostgres struct {
	db *sqlx.DB
}

func NewUpgradesPostgres(db *sqlx.DB) *UpgradesPostgres {
	return &UpgradesPostgres{db: db}
}

func (r *UpgradesPostgres) BuyUpgrade(userID int, upgradeID int) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, upgrade_id) VALUES ($1, $2)", userUpgradesTable)
	_, err := r.db.Exec(query, userID, upgradeID)
	return err
}

func (r *UpgradesPostgres) GetAllUpgrades() ([]model.Upgrade, error) {
	var upgrades []model.Upgrade
	query := fmt.Sprintf("SELECT * FROM %s", upgradesTable)
	if err := r.db.Select(&upgrades, query); err != nil {
		return nil, err
	}
	return upgrades, nil
}

func (r *UpgradesPostgres) GetUserUpgrades(userID int) ([]model.Upgrade, error) {
	var upgrades []model.Upgrade
	query := fmt.Sprintf("SELECT u.* FROM %s u JOIN %s uu ON u.id = uu.upgrade_id WHERE uu.user_id = $1", upgradesTable, userUpgradesTable)
	if err := r.db.Select(&upgrades, query, userID); err != nil {
		return nil, err
	}
	return upgrades, nil
}

func (r *UpgradesPostgres) CreateNewUpgrade(upgrade model.Upgrade) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, description, coins_price, stars_price, new_max_energy, new_one_tap_coins, new_passive_income) VALUES ($1, $2, $3, $4, $5, $6, $7)", upgradesTable)
	row := r.db.QueryRow(query, upgrade.Name, upgrade.Description, upgrade.CoinsPrice, upgrade.StarsPrice, upgrade.NewMaxEnergy, upgrade.NewOneTapCoins, upgrade.NewPassiveIncome)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
