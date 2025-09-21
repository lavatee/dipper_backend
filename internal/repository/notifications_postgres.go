package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dipper_backend/internal/model"
)

type NotificationsPostgres struct {
	db *sqlx.DB
}

func NewNotificationsPostgres(db *sqlx.DB) *NotificationsPostgres {
	return &NotificationsPostgres{
		db: db,
	}
}

func (r *NotificationsPostgres) CreateNotification(notification model.Notification) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, text) VALUES ($1, $2)", notificationsTable)
	_, err := r.db.Exec(query, notification.UserID, notification.Text)
	return err
}

func (r *NotificationsPostgres) GetUserNotifications(userID int) ([]model.Notification, error) {
	var nots []model.Notification
	query := fmt.Sprintf("SELECT id, text FROM %s WHERE user_id = $1", notificationsTable)
	if err := r.db.Select(&nots, query, userID); err != nil {
		return nil, err
	}
	return nots, nil
}
