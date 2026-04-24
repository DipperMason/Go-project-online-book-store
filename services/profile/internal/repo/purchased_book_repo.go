package repo

import (
	"context"
	"profile/internal/domain"
	"sort"
)

// PurchasedBookRepository интерфейс доступа к купленным книгам.
type PurchasedBookRepository interface {
	UpsertPurchasedBook(ctx context.Context, book *domain.PurchasedBook) error
	GetPurchasedBooks(ctx context.Context, userID int, limit int, offset int) ([]domain.PurchasedBook, error)
}

// MockPurchasedBookRepository in-memory реализация.
type MockPurchasedBookRepository struct {
	books map[int]map[int]domain.PurchasedBook
}

func NewMockPurchasedBookRepository() *MockPurchasedBookRepository {
	return &MockPurchasedBookRepository{books: make(map[int]map[int]domain.PurchasedBook)}
}

func (r *MockPurchasedBookRepository) UpsertPurchasedBook(ctx context.Context, book *domain.PurchasedBook) error {
	if _, ok := r.books[book.UserID]; !ok {
		r.books[book.UserID] = make(map[int]domain.PurchasedBook)
	}
	r.books[book.UserID][book.BookID] = *book
	return nil
}

func (r *MockPurchasedBookRepository) GetPurchasedBooks(ctx context.Context, userID int, limit int, offset int) ([]domain.PurchasedBook, error) {
	byBookID := r.books[userID]
	result := make([]domain.PurchasedBook, 0, len(byBookID))
	for _, book := range byBookID {
		result = append(result, book)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].PurchasedAt.After(result[j].PurchasedAt)
	})

	if offset >= len(result) {
		return []domain.PurchasedBook{}, nil
	}

	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], nil
}
