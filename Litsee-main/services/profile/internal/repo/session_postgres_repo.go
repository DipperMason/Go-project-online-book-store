package repo

import (
	"context"
	"database/sql"
	"profile/internal/domain"
)

// PostgresSessionRepository реализация SessionRepository для PostgreSQL
type PostgresSessionRepository struct {
	db *sql.DB
}

func NewPostgresSessionRepository(db *sql.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, token, ip_address, user_agent, is_active, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.Token,
		session.IPAddress,
		session.UserAgent,
		session.IsActive,
		session.CreatedAt,
		session.ExpiresAt,
	)

	if err != nil {
		return domain.ErrInternalServer
	}

	return nil
}

func (r *PostgresSessionRepository) GetActiveSessions(ctx context.Context, userID int) ([]domain.Session, error) {
	query := `
		SELECT id, user_id, token, ip_address, user_agent, created_at, expires_at, is_active
		FROM sessions
		WHERE user_id = $1 AND is_active = true AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	defer rows.Close()

	var sessions []domain.Session
	for rows.Next() {
		var session domain.Session
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Token,
			&session.IPAddress,
			&session.UserAgent,
			&session.CreatedAt,
			&session.ExpiresAt,
			&session.IsActive,
		)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *PostgresSessionRepository) TerminateSession(ctx context.Context, sessionID string) error {
	query := "UPDATE sessions SET is_active = false WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return domain.ErrInternalServer
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrInternalServer
	}

	if rowsAffected == 0 {
		return domain.ErrSessionNotFound
	}

	return nil
}

func (r *PostgresSessionRepository) TerminateAllSessions(ctx context.Context, userID int) error {
	query := "UPDATE sessions SET is_active = false WHERE user_id = $1"

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return domain.ErrInternalServer
	}

	return nil
}

func (r *PostgresSessionRepository) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	session := &domain.Session{}

	query := `
		SELECT id, user_id, token, ip_address, user_agent, created_at, expires_at, is_active
		FROM sessions
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.IPAddress,
		&session.UserAgent,
		&session.CreatedAt,
		&session.ExpiresAt,
		&session.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrSessionNotFound
	}
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	return session, nil
}
