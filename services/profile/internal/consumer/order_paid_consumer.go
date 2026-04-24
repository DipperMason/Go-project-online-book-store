package consumer

import (
	"context"
	"encoding/json"
	"log"
	"profile/internal/domain"
	"profile/internal/services"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// OrderPaidConsumer читает события orderPaid из Redpanda.
type OrderPaidConsumer struct {
	reader  *kafka.Reader
	service *services.ProfileService
}

func NewOrderPaidConsumer(brokers []string, topic string, groupID string, service *services.ProfileService) *OrderPaidConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &OrderPaidConsumer{reader: reader, service: service}
}

func (c *OrderPaidConsumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("orderPaid consumer read error: %v", err)
			continue
		}

		event, ok := parseOrderPaidEvent(msg.Value)
		if !ok {
			continue
		}

		if err = c.service.HandleOrderPaid(ctx, event); err != nil {
			log.Printf("orderPaid handle error: %v", err)
		}
	}
}

func (c *OrderPaidConsumer) Close() error {
	return c.reader.Close()
}

type rawEnvelope struct {
	Type      string          `json:"type"`
	EventType string          `json:"event_type"`
	Name      string          `json:"name"`
	Payload   json.RawMessage `json:"payload"`
	Data      json.RawMessage `json:"data"`
}

type rawOrderPaid struct {
	UserID  int           `json:"user_id"`
	OrderID string        `json:"order_id"`
	PaidAt  time.Time     `json:"paid_at"`
	Books   []rawBook     `json:"books"`
	BookIDs []int         `json:"book_ids"`
	BookID  int           `json:"book_id"`
	Items   []rawBook     `json:"items"`
	Meta    []interface{} `json:"meta"`
}

type rawBook struct {
	BookID   int    `json:"book_id"`
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	CoverURL string `json:"cover_url"`
}

func parseOrderPaidEvent(body []byte) (domain.OrderPaidEvent, bool) {
	var env rawEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return domain.OrderPaidEvent{}, false
	}

	eventType := normalizeEventType(firstNotEmpty(env.Type, env.EventType, env.Name))
	if eventType != "" && eventType != "orderpaid" {
		return domain.OrderPaidEvent{}, false
	}

	candidate := body
	if len(env.Payload) > 0 {
		candidate = env.Payload
	} else if len(env.Data) > 0 {
		candidate = env.Data
	}

	var raw rawOrderPaid
	if err := json.Unmarshal(candidate, &raw); err != nil {
		return domain.OrderPaidEvent{}, false
	}

	if raw.UserID <= 0 {
		return domain.OrderPaidEvent{}, false
	}

	books := make([]domain.OrderPaidBook, 0)
	for _, b := range append(raw.Books, raw.Items...) {
		bookID := b.BookID
		if bookID == 0 {
			bookID = b.ID
		}
		if bookID <= 0 {
			continue
		}
		books = append(books, domain.OrderPaidBook{
			BookID:   bookID,
			Title:    b.Title,
			Author:   b.Author,
			CoverURL: b.CoverURL,
		})
	}

	for _, id := range raw.BookIDs {
		if id <= 0 {
			continue
		}
		books = append(books, domain.OrderPaidBook{BookID: id})
	}

	if raw.BookID > 0 {
		books = append(books, domain.OrderPaidBook{BookID: raw.BookID})
	}

	if len(books) == 0 {
		return domain.OrderPaidEvent{}, false
	}

	paidAt := raw.PaidAt
	if paidAt.IsZero() {
		paidAt = time.Now().UTC()
	}

	return domain.OrderPaidEvent{
		UserID:  raw.UserID,
		OrderID: raw.OrderID,
		Books:   books,
		PaidAt:  paidAt,
	}, true
}

func firstNotEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func normalizeEventType(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	v = strings.ReplaceAll(v, "_", "")
	v = strings.ReplaceAll(v, "-", "")
	return v
}
