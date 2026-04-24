package domain

import "time"

// Session представляет активную сессию пользователя
type Session struct {
	ID        string
	UserID    int
	Token     string
	IPAddress string
	UserAgent string
	CreatedAt time.Time
	ExpiresAt time.Time
	IsActive  bool
}

// SessionInfo содержит информацию о сессии для отправки клиенту
type SessionInfo struct {
	ID        string    `json:"id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
}
