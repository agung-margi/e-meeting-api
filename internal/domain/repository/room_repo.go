package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
	"errors"
	"strconv"
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
	argIndex := 1

	// Filter berdasarkan name
	if name != "" {
		query += " AND name ILIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+name+"%")
		argIndex++
	}

	// Filter berdasarkan roomType
	if roomType > 0 {
		query += " AND room_type_id = $" + strconv.Itoa(argIndex)
		args = append(args, roomType)
		argIndex++
	}

	// Filter berdasarkan capacity
	if capacity > 0 {
		switch capacity {
		case 1:
			query += " AND capacity < $" + strconv.Itoa(argIndex)
			args = append(args, 10)
			argIndex++
		case 2:
			query += " AND capacity BETWEEN $" + strconv.Itoa(argIndex) + " AND $" + strconv.Itoa(argIndex+1)
			args = append(args, 11, 50)
			argIndex += 2
		case 3:
			query += " AND capacity BETWEEN $" + strconv.Itoa(argIndex) + " AND $" + strconv.Itoa(argIndex+1)
			args = append(args, 51, 100)
			argIndex += 2
		}
	}

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

func (r *roomRepo) DeleteRoom(ctx context.Context, id int) error {
	query := "DELETE FROM rooms WHERE id = $1"
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *reservationRepo) GetReservationsCountByRoomAndDate(ctx context.Context, roomID int, date string) (int, error) {
	query := `
			SELECT COUNT(id)
			FROM reservations 
			WHERE room_id = $1 
				AND DATE(start_time) = $2 
				AND (status = 'booked' OR status = 'paid')
	`

	var count int
	err := r.DB.QueryRowContext(ctx, query, roomID, date).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
