package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
	"errors"
)

type snackRepo struct {
	DB *sql.DB
}

func NewSnackRepository(db *sql.DB) SnackRepository {
	return &snackRepo{DB: db}
}

func (r *snackRepo) GetByID(ctx context.Context, id int) (*entity.Snack, error) {
	query := "SELECT id, name, category, price, created_at, updated_at FROM snacks WHERE id = $1"
	row := r.DB.QueryRowContext(ctx, query, id)
	snack := &entity.Snack{}

	err := row.Scan(
		&snack.ID,
		&snack.Name,
		&snack.Price,
		&snack.Category,
		&snack.CreatedAt,
		&snack.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("snack not found")
		}
		return nil, err
	}

	return snack, nil
}

func (r *snackRepo) GetAll(ctx context.Context) ([]entity.Snack, error) {
	query := "SELECT id, name, category, price, created_at, updated_at FROM snacks"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snacks []entity.Snack
	for rows.Next() {
		snack := entity.Snack{}
		if err := rows.Scan(
			&snack.ID,
			&snack.Name,
			&snack.Category,
			&snack.Price,
			&snack.CreatedAt,
			&snack.UpdatedAt,
		); err != nil {
			return nil, err
		}
		snacks = append(snacks, snack)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snacks, nil
}
