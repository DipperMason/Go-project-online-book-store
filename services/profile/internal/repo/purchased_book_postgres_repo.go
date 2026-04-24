package repo

import (
	"context"
	"database/sql"
	"profile/internal/domain"
)

// PostgresPurchasedBookRepository реализация PurchasedBookRepository для PostgreSQL.
type PostgresPurchasedBookRepository struct {
	db *sql.DB
}

func NewPostgresPurchasedBookRepository(db *sql.DB) *PostgresPurchasedBookRepository {
	return &PostgresPurchasedBookRepository{db: db}
}

func (r *PostgresPurchasedBookRepository) UpsertPurchasedBook(ctx context.Context, book *domain.PurchasedBook) error {
	query := `
		INSERT INTO purchased_books (user_id, order_id, book_id, title, author, cover_url, purchased_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, book_id)
		DO UPDATE SET
			order_id = EXCLUDED.order_id,
			title = EXCLUDED.title,
			author = EXCLUDED.author,
			cover_url = EXCLUDED.cover_url,
			purchased_at = EXCLUDED.purchased_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		book.UserID,
		book.OrderID,
		book.BookID,
		book.Title,
		book.Author,
		book.CoverURL,
		book.PurchasedAt,
	)
	if err != nil {
		return domain.ErrInternalServer
	}

	return nil
}

func (r *PostgresPurchasedBookRepository) GetPurchasedBooks(ctx context.Context, userID int, limit int, offset int) ([]domain.PurchasedBook, error) {
	query := `
		SELECT user_id, order_id, book_id, title, author, cover_url, purchased_at
		FROM purchased_books
		WHERE user_id = $1
		ORDER BY purchased_at DESC, book_id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	defer rows.Close()

	books := make([]domain.PurchasedBook, 0)
	for rows.Next() {
		var book domain.PurchasedBook
		if scanErr := rows.Scan(
			&book.UserID,
			&book.OrderID,
			&book.BookID,
			&book.Title,
			&book.Author,
			&book.CoverURL,
			&book.PurchasedAt,
		); scanErr != nil {
			continue
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.ErrInternalServer
	}

	return books, nil
}
