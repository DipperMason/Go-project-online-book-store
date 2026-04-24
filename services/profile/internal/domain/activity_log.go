package domain

import "time"

// ActivityLog представляет запись об активности пользователя
type ActivityLog struct {
	ID        int
	UserID    int
	Action    string // login, logout, profile_update, password_change
	IPAddress string
	UserAgent string
	Status    string // success, failed
	Details   string // дополнительная информация
	CreatedAt time.Time
}

// ActivityLogResponse используется для отправки истории активности
type ActivityLogResponse struct {
	ID        int       `json:"id"`
	Action    string    `json:"action"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Status    string    `json:"status"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}
