package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
	"errors"
)

type userRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{DB: db}
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	query := "SELECT id, username, email, password, is_admin, img_url, is_active, language, created_at, updated_at FROM users WHERE id = $1"
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
		&user.Language,
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
		INSERT INTO users (username, email, password, is_admin, img_url, is_active, language)
		VALUES ($1, $2, $3, $4, $5, $6, 'english')
		RETURNING id
	`
	err := r.DB.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.IsAdmin,
		user.ImgUrl,
		true,
	).Scan(&user.ID)
	return err
}

func (r *userRepo) UpdateUser(ctx context.Context, id int, user *entity.User) error {
	query := `
			UPDATE users
			SET username = $1, email = $2, password = $3, is_admin = $4, img_url = $5, is_active = $6, language = $7, updated_at = now()
			WHERE id = $8
	`
	_, err := r.DB.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.IsAdmin,
		user.ImgUrl,
		user.IsActive,
		user.Language,
		id,
	)
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
	query := "SELECT id, username, email, password, is_admin, img_url, language, is_active, created_at, updated_at FROM users WHERE username = $1"
	row := r.DB.QueryRowContext(ctx, query, username)
	user := &entity.User{}

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.ImgUrl,
		&user.Language,
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

func (r *userRepo) CheckUsernameExists(ctx context.Context, username string) (bool, error) {

	var user entity.User
	query := "SELECT id FROM users WHERE username = $1"

	err := r.DB.QueryRowContext(ctx, query, username).Scan(&user.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, id int, password string) error {
	query := "UPDATE users SET password = $1 WHERE id = $2"
	_, err := r.DB.ExecContext(ctx, query, password, id)
	return err
}
