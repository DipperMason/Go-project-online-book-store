package domain

import "time"

// PurchasedBook хранит книгу, купленную пользователем.
type PurchasedBook struct {
	UserID      int       `json:"user_id"`
	OrderID     string    `json:"order_id"`
	BookID      int       `json:"book_id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	CoverURL    string    `json:"cover_url"`
	PurchasedAt time.Time `json:"purchased_at"`
}

// OrderPaidEvent - событие успешной оплаты заказа.
type OrderPaidEvent struct {
	UserID  int
	OrderID string
	Books   []OrderPaidBook
	PaidAt  time.Time
}

// OrderPaidBook - книга из события orderPaid.
type OrderPaidBook struct {
	BookID   int
	Title    string
	Author   string
	CoverURL string
}
