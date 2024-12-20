package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
	"errors"
	"time"
)

type userRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{DB: db}
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	query := "SELECT id, username, email, password, is_admin, img_url, is_active, created_at, updated_at FROM users WHERE id = $1"
	row := r.DB.QueryRowContext(ctx, query, id)
	user := &entity.User{}

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.ImgUrl,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepo) SaveUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (username, email, password, is_admin, img_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.DB.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.IsAdmin,
		user.ImgUrl,
		user.IsActive,
	).Scan(&user.ID)
	return err
}

func (r *userRepo) UpdateUser(ctx context.Context, id int, user *entity.User) error {
	query := `
			UPDATE users
			SET password = $1, img_url = $2, is_active = $3, updated_at = $4
			WHERE id = $5
	`
	_, err := r.DB.ExecContext(ctx, query,
		user.Password,
		user.ImgUrl,
		user.IsActive,
		time.Now(), // Use current timestamp for updated_at
		id,
	)
	return err
}

func (r *userRepo) DeleteUser(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

func (r *userRepo) CheckEmailExists(ctx context.Context, email string) (bool, error) {

	var user entity.User
	query := "SELECT id FROM users WHERE email = $1"

	err := r.DB.QueryRowContext(ctx, query, email).Scan(&user.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := "SELECT id, username, email, password, is_admin, img_url, is_active, created_at, updated_at FROM users WHERE username = $1"
	row := r.DB.QueryRowContext(ctx, query, username)
	user := &entity.User{}

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.ImgUrl,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
