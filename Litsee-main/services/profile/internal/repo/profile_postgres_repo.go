package repo

import (
	"context"
	"database/sql"
	"profile/internal/domain"
	"time"
)

// PostgresProfileRepository реализация ProfileRepository для PostgreSQL.
type PostgresProfileRepository struct {
	db *sql.DB
}

func NewPostgresProfileRepository(db *sql.DB) *PostgresProfileRepository {
	return &PostgresProfileRepository{db: db}
}

func (r *PostgresProfileRepository) GetProfile(ctx context.Context, userID int) (*domain.UserProfile, error) {
	profile := &domain.UserProfile{}

	query := `
		SELECT id, first_name, last_name, avatar, bio, updated_at, created_at
		FROM user_profiles
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Avatar,
		&profile.Bio,
		&profile.UpdatedAt,
		&profile.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return &domain.UserProfile{ID: userID}, nil
	}
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	return profile, nil
}

func (r *PostgresProfileRepository) UpdateProfile(ctx context.Context, userID int, profile *domain.UpdateProfileRequest) error {
	query := `
		INSERT INTO user_profiles (id, first_name, last_name, avatar, bio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id)
		DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			avatar = EXCLUDED.avatar,
			bio = EXCLUDED.bio,
			updated_at = EXCLUDED.updated_at
	`

	now := time.Now().UTC()
	_, err := r.db.ExecContext(
		ctx,
		query,
		userID,
		profile.FirstName,
		profile.LastName,
		profile.Avatar,
		profile.Bio,
		now,
		now,
	)
	if err != nil {
		return domain.ErrInternalServer
	}

	return nil
}
