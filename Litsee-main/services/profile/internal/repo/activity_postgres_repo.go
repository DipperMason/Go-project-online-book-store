package repo

import (
	"context"
	"database/sql"
	"profile/internal/domain"
)

// PostgresActivityRepository реализация ActivityRepository для PostgreSQL
type PostgresActivityRepository struct {
	db *sql.DB
}

func NewPostgresActivityRepository(db *sql.DB) *PostgresActivityRepository {
	return &PostgresActivityRepository{db: db}
}

func (r *PostgresActivityRepository) LogActivity(ctx context.Context, log *domain.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (user_id, action, ip_address, user_agent, status, details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		log.UserID,
		log.Action,
		log.IPAddress,
		log.UserAgent,
		log.Status,
		log.Details,
		log.CreatedAt,
	)
	
	if err != nil {
		return domain.ErrInternalServer
	}
	
	return nil
}

func (r *PostgresActivityRepository) GetActivityHistory(ctx context.Context, userID int, limit int, offset int) ([]domain.ActivityLog, error) {
	query := `
		SELECT id, user_id, action, ip_address, user_agent, status, details, created_at
		FROM activity_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	defer rows.Close()
	
	var logs []domain.ActivityLog
	for rows.Next() {
		var log domain.ActivityLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.IPAddress,
			&log.UserAgent,
			&log.Status,
			&log.Details,
			&log.CreatedAt,
		)
		if err != nil {
			continue
		}
		logs = append(logs, log)
	}
	
	return logs, nil
}
