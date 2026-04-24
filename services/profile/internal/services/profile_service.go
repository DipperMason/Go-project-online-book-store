package services

import (
	"context"
	"profile/internal/domain"
	"profile/internal/repo"
	"time"
)

// ProfileService сервис для работы с профилем пользователя.
type ProfileService struct {
	profileRepo       repo.ProfileRepository
	activityRepo      repo.ActivityRepository
	purchasedBookRepo repo.PurchasedBookRepository
}

func NewProfileService(
	profileRepo repo.ProfileRepository,
	activityRepo repo.ActivityRepository,
	purchasedBookRepo repo.PurchasedBookRepository,
) *ProfileService {
	return &ProfileService{
		profileRepo:       profileRepo,
		activityRepo:      activityRepo,
		purchasedBookRepo: purchasedBookRepo,
	}
}

// GetProfile получает профиль пользователя.
func (s *ProfileService) GetProfile(ctx context.Context, userID int) (*domain.UserProfile, error) {
	profile, err := s.profileRepo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdateProfile обновляет профиль пользователя.
func (s *ProfileService) UpdateProfile(ctx context.Context, userID int, req *domain.UpdateProfileRequest) error {
	err := s.profileRepo.UpdateProfile(ctx, userID, req)
	if err != nil {
		return err
	}

	_ = s.activityRepo.LogActivity(ctx, &domain.ActivityLog{
		UserID:    userID,
		Action:    "profile_update",
		Status:    "success",
		CreatedAt: time.Now().UTC(),
	})

	return nil
}

// GetPurchasedBooks возвращает библиотеку купленных книг пользователя.
func (s *ProfileService) GetPurchasedBooks(ctx context.Context, userID int, limit int, offset int) ([]domain.PurchasedBook, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	books, err := s.purchasedBookRepo.GetPurchasedBooks(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	if books == nil {
		return []domain.PurchasedBook{}, nil
	}
	return books, nil
}

// HandleOrderPaid сохраняет книги из события orderPaid.
func (s *ProfileService) HandleOrderPaid(ctx context.Context, event domain.OrderPaidEvent) error {
	paidAt := event.PaidAt
	if paidAt.IsZero() {
		paidAt = time.Now().UTC()
	}

	for _, book := range event.Books {
		err := s.purchasedBookRepo.UpsertPurchasedBook(ctx, &domain.PurchasedBook{
			UserID:      event.UserID,
			OrderID:     event.OrderID,
			BookID:      book.BookID,
			Title:       book.Title,
			Author:      book.Author,
			CoverURL:    book.CoverURL,
			PurchasedAt: paidAt,
		})
		if err != nil {
			return err
		}
	}

	_ = s.activityRepo.LogActivity(ctx, &domain.ActivityLog{
		UserID:    event.UserID,
		Action:    "order_paid",
		Status:    "success",
		Details:   event.OrderID,
		CreatedAt: time.Now().UTC(),
	})

	return nil
}

// GetActivityHistory получает историю активности пользователя.
func (s *ProfileService) GetActivityHistory(ctx context.Context, userID int, limit int, offset int) ([]domain.ActivityLogResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := s.activityRepo.GetActivityHistory(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.ActivityLogResponse, 0, len(logs))
	for _, log := range logs {
		responses = append(responses, domain.ActivityLogResponse{
			ID:        log.ID,
			Action:    log.Action,
			IPAddress: log.IPAddress,
			UserAgent: log.UserAgent,
			Status:    log.Status,
			Details:   log.Details,
			CreatedAt: log.CreatedAt,
		})
	}
	return responses, nil
}
