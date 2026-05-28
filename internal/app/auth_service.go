package app

import (
	"context"
	"docintel/internal/domain"
	"docintel/pkg/auth"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	InvalidCredentialsErr = "invalid email or password"
)

type AuthService struct {
	userRepo domain.UserRepository
	cache    domain.Cache
}

func NewAuthService(userRepo domain.UserRepository, cache domain.Cache) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cache:    cache,
	}
}

func (s *AuthService) Register(ctx context.Context, email, name, password string) error {
	passwordHash := hashPassword(password)
	newUser := &domain.User{
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
	}

	return s.userRepo.Create(ctx, newUser)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(InvalidCredentialsErr) // Invalid credentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New(InvalidCredentialsErr) // Invalid credentials
	}

	SaveUserSession(ctx, s.cache, user.ID.String())

	return user, nil
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func SaveUserSession(ctx context.Context, cache domain.Cache, userID string) (string, error) {
	token, err := auth.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	sessionID := "session:" + token
	err = cache.Set(ctx, sessionID, userID, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func ValidateUserSession(ctx context.Context, cache domain.Cache, sessionID string) (string, error) {
	userID, err := cache.Get(ctx, sessionID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func InvalidateUserSession(ctx context.Context, cache domain.Cache, sessionID string) error {
	return cache.Delete(ctx, sessionID)
}
