package repo

import (
	"context"
	"profile/internal/domain"
	"time"
)

// SessionRepository интерфейс для работы с сессиями
type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetActiveSessions(ctx context.Context, userID int) ([]domain.Session, error)
	TerminateSession(ctx context.Context, sessionID string) error
	TerminateAllSessions(ctx context.Context, userID int) error
	GetSession(ctx context.Context, sessionID string) (*domain.Session, error)
}

// MockSessionRepository использует в памяти хранилище
type MockSessionRepository struct {
	sessions map[string]*domain.Session
}

func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions: make(map[string]*domain.Session),
	}
}

func (r *MockSessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	r.sessions[session.ID] = session
	return nil
}

func (r *MockSessionRepository) GetActiveSessions(ctx context.Context, userID int) ([]domain.Session, error) {
	var sessions []domain.Session
	for _, session := range r.sessions {
		if session.UserID == userID && session.IsActive && session.ExpiresAt.After(time.Now()) {
			sessions = append(sessions, *session)
		}
	}
	return sessions, nil
}

func (r *MockSessionRepository) TerminateSession(ctx context.Context, sessionID string) error {
	if session, ok := r.sessions[sessionID]; ok {
		session.IsActive = false
		return nil
	}
	return domain.ErrSessionNotFound
}

func (r *MockSessionRepository) TerminateAllSessions(ctx context.Context, userID int) error {
	for _, session := range r.sessions {
		if session.UserID == userID {
			session.IsActive = false
		}
	}
	return nil
}

func (r *MockSessionRepository) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	if session, ok := r.sessions[sessionID]; ok {
		return session, nil
	}
	return nil, domain.ErrSessionNotFound
}
