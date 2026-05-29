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

	err := s.userRepo.CreateUser(ctx, newUser)
	return err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New(InvalidCredentialsErr) // Invalid credentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New(InvalidCredentialsErr) // Invalid credentials
	}

	token, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		return nil, "", err
	}
	if err := s.cache.Set(ctx, "session:"+token, user.ID.String(), 24*time.Hour); err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) SaveUserSession(ctx context.Context, userID string) (string, error) {
	token, err := auth.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	sessionID := "session:" + token
	err = s.cache.Set(ctx, sessionID, userID, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *AuthService) ValidateUserSession(ctx context.Context, sessionID string) (string, error) {
	userID, err := s.cache.Get(ctx, sessionID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *AuthService) InvalidateUserSession(ctx context.Context, sessionID string) error {
	return s.cache.Delete(ctx, sessionID)
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
