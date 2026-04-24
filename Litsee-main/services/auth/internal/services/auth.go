package services

import (
	"auth/internal/domain"
	"jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	GetUser(id int) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	CreateUser(email, password string) (*domain.User, error)
}

type AuthService struct {
	userRepo   UserRepo
	jwtService jwt.Service
}

func NewAuthService(userRepo UserRepo, jwtService jwt.Service) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (a *AuthService) Register(email string, password string) error {
	_, err := a.userRepo.CreateUser(email, password)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Authenticate(email string, password string) (jwt.TokenPair, error) {
	userInDB, err := a.userRepo.GetUserByEmail(email)
	if err != nil {
		return jwt.TokenPair{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInDB.PasswordHash), []byte(password))
	if err != nil {
		return jwt.TokenPair{}, domain.ErrInvalidCredentials
	}

	tokenPair, err := a.jwtService.GenerateToken(jwt.UserClaims{
		UserID:  userInDB.ID,
		Email:   userInDB.Email,
		IsAdmin: userInDB.IsAdmin,
	})

	if err != nil {
		return jwt.TokenPair{}, err
	}

	return tokenPair, nil
}
