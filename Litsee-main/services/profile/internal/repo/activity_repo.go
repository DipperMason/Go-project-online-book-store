package repo

import (
	"context"
	"profile/internal/domain"
)

// ActivityRepository интерфейс для работы с логами активности
type ActivityRepository interface {
	LogActivity(ctx context.Context, log *domain.ActivityLog) error
	GetActivityHistory(ctx context.Context, userID int, limit int, offset int) ([]domain.ActivityLog, error)
}

// MockActivityRepository использует в памяти хранилище
type MockActivityRepository struct {
	logs   []domain.ActivityLog
	nextID int
}

func NewMockActivityRepository() *MockActivityRepository {
	return &MockActivityRepository{
		logs:   make([]domain.ActivityLog, 0),
		nextID: 1,
	}
}

func (r *MockActivityRepository) LogActivity(ctx context.Context, log *domain.ActivityLog) error {
	log.ID = r.nextID
	r.nextID++
	r.logs = append(r.logs, *log)
	return nil
}

func (r *MockActivityRepository) GetActivityHistory(ctx context.Context, userID int, limit int, offset int) ([]domain.ActivityLog, error) {
	var userLogs []domain.ActivityLog

	// Собираем логи пользователя в обратном порядке (новые первыми)
	for i := len(r.logs) - 1; i >= 0; i-- {
		if r.logs[i].UserID == userID {
			userLogs = append(userLogs, r.logs[i])
		}
	}

	// Применяем пагинацию
	end := offset + limit
	if end > len(userLogs) {
		end = len(userLogs)
	}
	if offset > len(userLogs) {
		offset = len(userLogs)
	}

	return userLogs[offset:end], nil
}
