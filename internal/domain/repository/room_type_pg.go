package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
	"errors"
)

type roomTypeRepo struct {
	DB *sql.DB
}

func NewRoomtypeRepository(db *sql.DB) RoomtypeRepository {
	return &roomTypeRepo{DB: db}
}

func (r *roomTypeRepo) GetAll(ctx context.Context) ([]entity.RoomtypeRepository, error) {
	query := "SELECT id, name, created_at, updated_at FROM room_types"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roomTypes []entity.RoomtypeRepository
	for rows.Next() {
		roomType := entity.RoomtypeRepository{}
		if err := rows.Scan(
			&roomType.ID,
			&roomType.Name,
			&roomType.CreatedAt,
			&roomType.UpdatedAt,
		); err != nil {
			return nil, err
		}
		roomTypes = append(roomTypes, roomType)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roomTypes, nil
}

func (r *roomTypeRepo) GetByID(ctx context.Context, id int) (*entity.RoomtypeRepository, error) {
	query := "SELECT id, name, created_at, updated_at FROM room_types WHERE id = $1"
	row := r.DB.QueryRowContext(ctx, query, id)
	roomType := &entity.RoomtypeRepository{}

	err := row.Scan(
		&roomType.ID,
		&roomType.Name,
		&roomType.CreatedAt,
		&roomType.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("room type not found")
		}
		return nil, err
	}

	return roomType, nil
}

func (r *roomTypeRepo) GetByName(ctx context.Context, name string) (*entity.RoomtypeRepository, error) {
	query := "SELECT id, name, created_at, updated_at FROM room_types WHERE name = $1"
	row := r.DB.QueryRowContext(ctx, query, name)
	roomType := &entity.RoomtypeRepository{}

	err := row.Scan(
		&roomType.ID,
		&roomType.Name,
		&roomType.CreatedAt,
		&roomType.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("room type not found")
		}
		return nil, err
	}

	return roomType, nil
}

func (r *roomTypeRepo) SaveRoomType(ctx context.Context, roomType *entity.RoomtypeRepository) error {
	query := "INSERT INTO room_types (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id"
	row := r.DB.QueryRowContext(ctx, query, roomType.Name, roomType.CreatedAt, roomType.UpdatedAt)
	if err := row.Scan(&roomType.ID); err != nil {
		return err
	}
	return nil
}
