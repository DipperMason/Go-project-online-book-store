package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"profile/internal/config"
	"profile/internal/consumer"
	"profile/internal/repo"
	"profile/internal/services"
	httphandler "profile/internal/transport/http"
	"time"

	"jwt"

	_ "github.com/lib/pq"
)

// App представляет приложение профиля.
type App struct {
	handler  http.Handler
	db       *sql.DB
	consumer *consumer.OrderPaidConsumer
	cancel   context.CancelFunc
}

func New() (*App, error) {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	pingCtx, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()
	if err = db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	profileRepo := repo.NewPostgresProfileRepository(db)
	activityRepo := repo.NewPostgresActivityRepository(db)
	purchasedBookRepo := repo.NewPostgresPurchasedBookRepository(db)

	profileService := services.NewProfileService(profileRepo, activityRepo, purchasedBookRepo)
	jwtService := jwt.NewService(jwt.Config{Secret: cfg.JWTSecret})
	handler := httphandler.New(profileService, jwtService)

	application := &App{
		handler: handler.Router(),
		db:      db,
	}

	if cfg.RedpandaOrderPaidTopic != "" && len(cfg.RedpandaBrokers) > 0 {
		ctx, cancel := context.WithCancel(context.Background())
		orderPaidConsumer := consumer.NewOrderPaidConsumer(
			cfg.RedpandaBrokers,
			cfg.RedpandaOrderPaidTopic,
			cfg.RedpandaConsumerGroupID,
			profileService,
		)
		application.consumer = orderPaidConsumer
		application.cancel = cancel
		go orderPaidConsumer.Run(ctx)
	}

	return application, nil
}

func (a *App) GetHandler() http.Handler {
	return a.handler
}

func (a *App) Close() error {
	if a.cancel != nil {
		a.cancel()
	}
	if a.consumer != nil {
		_ = a.consumer.Close()
	}
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}
