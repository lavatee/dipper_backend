package model

type Notification struct {
	ID     int    `json:"id" db:"id"`
	UserID int    `json:"user_id" db:"user_id"`
	Text   string `json:"text" db:"text"`
}
