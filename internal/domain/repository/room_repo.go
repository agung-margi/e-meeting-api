package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
	"errors"
	"fmt"
)

type roomRepo struct {
	DB *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepo{DB: db}
}

func (r *roomRepo) GetAll(ctx context.Context, name string, roomType int, capacity int) ([]entity.Room, error) {
	query := "SELECT id, name, room_type_id, price, capacity, img_url, created_at, updated_at FROM rooms WHERE 1=1"

	var args []interface{}
	argIndex := 1 // Counter untuk PostgreSQL placeholder ($1, $2, dll.)

	// Filter berdasarkan name
	if name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+name+"%")
		argIndex++
	}

	// Filter berdasarkan roomType
	if roomType > 0 {
		query += fmt.Sprintf(" AND room_type_id = $%d", argIndex)
		args = append(args, roomType)
		argIndex++
	}

	// Filter berdasarkan capacity
	if capacity > 0 {
		switch capacity {
		case 1:
			query += fmt.Sprintf(" AND capacity < $%d", argIndex)
			args = append(args, 10)
			argIndex++
		case 2:
			query += fmt.Sprintf(" AND capacity BETWEEN $%d AND $%d", argIndex, argIndex+1)
			args = append(args, 11, 50)
			argIndex += 2
		case 3:
			query += fmt.Sprintf(" AND capacity BETWEEN $%d AND $%d", argIndex, argIndex+1)
			args = append(args, 51, 100)
			argIndex += 2
		}
	}

	// Log query untuk debugging
	fmt.Printf("Executing query: %s, with args: %v\n", query, args)

	// Eksekusi query
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parsing hasil query
	var rooms []entity.Room
	for rows.Next() {
		room := entity.Room{}
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.RoomType,
			&room.Price,
			&room.Capacity,
			&room.ImgUrl,
			&room.CreatedAt,
			&room.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *roomRepo) GetByID(ctx context.Context, id int) (*entity.RoomWithType, error) {
	query := `
		SELECT 
			r.id, r.name, rt.name AS room_type, 
			r.price, r.capacity, r.img_url, 
			r.created_at, r.updated_at
		FROM 
			rooms r
		INNER JOIN 
			room_types rt ON r.room_type_id = rt.id
		WHERE 
			r.id = $1
	`
	row := r.DB.QueryRowContext(ctx, query, id)

	room := &entity.RoomWithType{}

	err := row.Scan(
		&room.ID,
		&room.Name,
		&room.RoomType, // Room type name instead of ID
		&room.Price,
		&room.Capacity,
		&room.ImgUrl,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}
	return room, nil
}

func (r *roomRepo) SaveRoom(ctx context.Context, room *entity.Room) error {
	query := "INSERT INTO rooms (name, room_type_id, price, capacity, img_url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	row := r.DB.QueryRowContext(ctx, query, room.Name, room.RoomType, room.Price, room.Capacity, room.ImgUrl, room.CreatedAt, room.UpdatedAt)
	if err := row.Scan(&room.ID); err != nil {
		return err
	}
	return nil
}

func (r *roomRepo) UpdateRoom(ctx context.Context, id int, room *entity.Room) error {
	query := `
		UPDATE rooms
		SET name = $1, room_type_id = $2, price = $3, capacity = $4, img_url = $5, updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.DB.ExecContext(ctx, query, room.Name, room.RoomType, room.Price, room.Capacity, room.ImgUrl, room.ID)
	if err != nil {
		return err
	}
	return nil
}
