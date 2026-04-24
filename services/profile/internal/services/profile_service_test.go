package services

import (
	"context"
	"profile/internal/domain"
	"profile/internal/repo"
	"testing"
	"time"
)

func TestGetProfile(t *testing.T) {
	profileRepo := repo.NewMockProfileRepository()
	activityRepo := repo.NewMockActivityRepository()
	booksRepo := repo.NewMockPurchasedBookRepository()
	service := NewProfileService(profileRepo, activityRepo, booksRepo)

	profile, err := service.GetProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if profile == nil || profile.ID != 1 {
		t.Fatalf("expected profile with ID=1, got %+v", profile)
	}
}

func TestUpdateProfile(t *testing.T) {
	profileRepo := repo.NewMockProfileRepository()
	activityRepo := repo.NewMockActivityRepository()
	booksRepo := repo.NewMockPurchasedBookRepository()
	service := NewProfileService(profileRepo, activityRepo, booksRepo)

	userID := 1
	err := service.UpdateProfile(context.Background(), userID, &domain.UpdateProfileRequest{
		FirstName: "John",
		LastName:  "Doe",
		Bio:       "Test bio",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	profile, _ := service.GetProfile(context.Background(), userID)
	if profile.FirstName != "John" {
		t.Fatalf("expected FirstName John, got %s", profile.FirstName)
	}
}

func TestHandleOrderPaidAndGetPurchasedBooks(t *testing.T) {
	profileRepo := repo.NewMockProfileRepository()
	activityRepo := repo.NewMockActivityRepository()
	booksRepo := repo.NewMockPurchasedBookRepository()
	service := NewProfileService(profileRepo, activityRepo, booksRepo)

	err := service.HandleOrderPaid(context.Background(), domain.OrderPaidEvent{
		UserID:  7,
		OrderID: "order-42",
		PaidAt:  time.Now().UTC(),
		Books: []domain.OrderPaidBook{
			{BookID: 11, Title: "DDD", Author: "Evans"},
			{BookID: 12, Title: "Clean Code", Author: "Martin"},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	books, err := service.GetPurchasedBooks(context.Background(), 7, 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
}

func TestGetActivityHistory(t *testing.T) {
	profileRepo := repo.NewMockProfileRepository()
	activityRepo := repo.NewMockActivityRepository()
	booksRepo := repo.NewMockPurchasedBookRepository()
	service := NewProfileService(profileRepo, activityRepo, booksRepo)

	userID := 1
	_ = activityRepo.LogActivity(context.Background(), &domain.ActivityLog{UserID: userID, Action: "login", Status: "success"})
	_ = activityRepo.LogActivity(context.Background(), &domain.ActivityLog{UserID: userID, Action: "profile_update", Status: "success"})

	history, err := service.GetActivityHistory(context.Background(), userID, 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("expected 2 activities, got %d", len(history))
	}
}
