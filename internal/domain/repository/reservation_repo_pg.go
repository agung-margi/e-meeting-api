package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
	"errors"
)

type reservationRepo struct {
	DB *sql.DB
}

func NewReservationRepository(db *sql.DB) ReservationRepository {
	return &reservationRepo{DB: db}
}

func (r *reservationRepo) GetByID(ctx context.Context, id int) (*entity.Reservation, error) {
	query := "SELECT * FROM reservations WHERE id = $1"
	row := r.DB.QueryRowContext(ctx, query, id)
	roomReservation := &entity.Reservation{}
	err := row.Scan(
		&roomReservation.ID,
		&roomReservation.UserID,
		&roomReservation.RoomID,
		&roomReservation.StartTime,
		&roomReservation.EndTime,
		&roomReservation.BookingDate,
		&roomReservation.RoomPrice,
		&roomReservation.SnackPrice,
		&roomReservation.TotalPrice,
		&roomReservation.Status,
		&roomReservation.CreatedAt,
		&roomReservation.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return roomReservation, nil
}

func (r *reservationRepo) SaveReservation(ctx context.Context, rsv *entity.Reservation) error {
	query := "INSERT INTO reservations (user_id, room_id, start_time, end_time, booking_date, room_price, snack_price, total_price, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	err := r.DB.QueryRowContext(ctx, query,
		rsv.UserID,
		rsv.RoomID,
		rsv.StartTime,
		rsv.EndTime,
		rsv.BookingDate,
		rsv.RoomPrice,
		rsv.SnackPrice,
		rsv.TotalPrice,
		rsv.Status,
	).Scan(&rsv.ID)

	return err
}

// func (r *reservationRepo) SaveReservationDetail(ctx context.Context, rsvd *entity.ReservationDetail) error {
// 	query := "INSERT INTO reservation_details (resevation_id, name, phone, company, snack_id, participants) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
// 	err := r.DB.QueryRowContext(ctx, query,
// 		rsvd.ReservationID,
// 		rsvd.Name,
// 		rsvd.Phone,
// 		rsvd.Company,
// 		rsvd.SnackID,
// 		rsvd.Participants,
// 	).Scan(&rsvd.ID)
// 	return err
// }

func (r *reservationRepo) CheckAvailability(ctx context.Context, roomId int, startTime string, endTime string) (bool, error) {
	query := `
        SELECT 1
        FROM reservations
        WHERE room_id = $1
          AND (start_time < $3 AND end_time > $2)
        LIMIT 1
    `
	var exists int
	err := r.DB.QueryRowContext(ctx, query, roomId, startTime, endTime).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, err
}
func (r *reservationRepo) GetReservationsByRoomAndDate(ctx context.Context, roomID int, date string) ([]*entity.Reservation, error) {
	query := `
			SELECT id, user_id, room_id, start_time, end_time, booking_date, room_price, snack_price, total_price, status, created_at, updated_at
			FROM reservations
			WHERE room_id = $1
				AND DATE(start_time) = $2
			ORDER BY start_time ASC
	`
	rows, err := r.DB.QueryContext(ctx, query, roomID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*entity.Reservation
	for rows.Next() {
		var reservation entity.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.RoomID,
			&reservation.StartTime,
			&reservation.EndTime,
			&reservation.BookingDate,
			&reservation.RoomPrice,
			&reservation.SnackPrice,
			&reservation.TotalPrice,
			&reservation.Status,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}
