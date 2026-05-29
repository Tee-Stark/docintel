package postgres

import (
	"context"
	"database/sql"
	"docintel/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, email, name, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	user_id, uuid_err := uuid.NewV7()
	if uuid_err != nil {
		return uuid_err
	}

	_, exec_err := r.db.ExecContext(ctx, query, user_id, user.Email, user.Name, user.PasswordHash)
	return exec_err
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, name, password_hash, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindAllUsers(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, email, name, created_at, updated_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	panic("not implemented")
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}
