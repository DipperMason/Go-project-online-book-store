package repo

import (
	"context"
	"profile/internal/domain"
)

// ProfileRepository интерфейс для работы с профилями.
type ProfileRepository interface {
	GetProfile(ctx context.Context, userID int) (*domain.UserProfile, error)
	UpdateProfile(ctx context.Context, userID int, profile *domain.UpdateProfileRequest) error
}

// MockProfileRepository использует in-memory хранилище.
type MockProfileRepository struct {
	profiles map[int]*domain.UserProfile
}

func NewMockProfileRepository() *MockProfileRepository {
	return &MockProfileRepository{
		profiles: make(map[int]*domain.UserProfile),
	}
}

func (r *MockProfileRepository) GetProfile(ctx context.Context, userID int) (*domain.UserProfile, error) {
	if profile, ok := r.profiles[userID]; ok {
		return profile, nil
	}
	return &domain.UserProfile{ID: userID}, nil
}

func (r *MockProfileRepository) UpdateProfile(ctx context.Context, userID int, profile *domain.UpdateProfileRequest) error {
	existing, _ := r.GetProfile(ctx, userID)
	existing.FirstName = profile.FirstName
	existing.LastName = profile.LastName
	existing.Avatar = profile.Avatar
	existing.Bio = profile.Bio
	r.profiles[userID] = existing
	return nil
}
