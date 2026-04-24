package repo

import (
	"auth/internal/domain"
	"math/rand"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	fakeDb map[string]domain.User
}

func New() *UserRepo {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)

	return &UserRepo{
		fakeDb: map[string]domain.User{
			"reader@litsee.local": {
				ID:           1,
				Email:        "reader@litsee.local",
				PasswordHash: string(hashedPassword),
				IsAdmin:      false,
			},
		},
	}
}

func (u *UserRepo) GetUser(id int) (*domain.User, error) {
	sq.Select("*").From("users").Where(sq.Eq{"id": id})
	return &domain.User{
		ID:           1,
		Email:        "test@mai.ru",
		PasswordHash: "dsgfdsfg",
		IsAdmin:      false,
	}, nil
}

func (u *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	//sq.Select("*").From("users").Where(sq.Eq{"email": email})
	//return &domain.User{
	//	ID:           1,
	//	Email:        "test@mai.ru",
	//	PasswordHash: "dsgfdsfg",
	//	IsAdmin:      false,
	//}, nil

	if user, ok := u.fakeDb[email]; ok {
		return &user, nil
	}

	return nil, domain.UserNotFound
}

func (u *UserRepo) CreateUser(email string, password string) (*domain.User, error) {
	if _, ok := u.fakeDb[email]; ok {
		return nil, domain.UserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		ID:           int(rand.Int31()),
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsAdmin:      false,
	}

	u.fakeDb[email] = user

	return &user, nil
}
