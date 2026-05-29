package domain

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindAllUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, doc *Document) error
	// FindDocumentByID(ctx context.Context, id uuid.UUID) (*Document, error)
}

type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type AuthService interface {
	Register(ctx context.Context, email, name, password string) error
	Login(ctx context.Context, email, password string) (*User, error)
	SaveUserSession(ctx context.Context, userID string) (string, error)
	ValidateUserSession(ctx context.Context, sessionID string) (string, error)
	InvalidateUserSession(ctx context.Context, sessionID string) error
}

type DocumentService interface {
	UploadDocument(ctx context.Context, file multipart.File, document *Document) error
}

type Repository interface {
	UserRepository
	DocumentRepository
}
