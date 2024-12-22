package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/presenter/model"
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

func (r *reservationRepo) GetSnackPriceByID(ctx context.Context, snackID int) (int, error) {
	var price int
	query := "SELECT price FROM snacks WHERE id = $1"
	err := r.DB.QueryRowContext(ctx, query, snackID).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return price, nil
}

func (r *reservationRepo) GetRoomPriceByID(ctx context.Context, roomID int) (int, error) {
	var price int
	query := "SELECT price FROM rooms WHERE id = $1"
	err := r.DB.QueryRowContext(ctx, query, roomID).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return price, nil
}
func (r *reservationRepo) SaveReservation(ctx context.Context, reservation *entity.Reservation, details *entity.ReservationDetails) error {
	query := "INSERT INTO reservations (user_id, room_id, start_time, end_time, booking_date, room_price, snack_price, total_price, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now()) RETURNING id"
	row := r.DB.QueryRowContext(ctx, query,
		reservation.UserID,
		reservation.RoomID,
		reservation.StartTime,
		reservation.EndTime,
		reservation.BookingDate,
		reservation.RoomPrice,
		reservation.SnackPrice,
		reservation.TotalPrice,
		"booked",
	)
	if err := row.Scan(&reservation.ID); err != nil {
		return err
	}

	details.ReservationID = reservation.ID

	querydetails := "INSERT INTO reservation_details (reservation_id, name, phone, company, snack_id, participants, notes) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	detail := r.DB.QueryRowContext(ctx, querydetails,
		details.ReservationID,
		details.Name,
		details.Phone,
		details.Company,
		details.SnackID,
		details.Participants,
		details.Notes)
	if err := detail.Scan(&details.ID); err != nil {
		return err
	}

	return nil
}

func (r *reservationRepo) GetReservationDetails(ctx context.Context, reservationID int) ([]model.ReservationDetailsResponse, error) {
	query := `
			SELECT
					rd.reservation_id,
					rd.name,
					rd.phone,
					rd.company,
					rd.snack_id,
					rd.participants,
					rd.notes
			FROM
					reservation_details rd
			WHERE
					rd.reservation_id = $1
	`

	rows, err := r.DB.QueryContext(ctx, query, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []model.ReservationDetailsResponse
	for rows.Next() {
		var detail model.ReservationDetailsResponse
		var snackID sql.NullInt64
		if err := rows.Scan(
			&detail.ReservationID,
			&detail.Name,
			&detail.Phone,
			&detail.Company,
			&snackID,
			&detail.Participants,
			&detail.Notes,
		); err != nil {
			return nil, err
		}

		if snackID.Valid {
			detail.SnackID = int(snackID.Int64)
		} else {
			detail.SnackID = 0
		}
		details = append(details, detail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return details, nil
}

func (r *reservationRepo) GetAll(ctx context.Context) ([]*entity.Reservation, error) {
	query := "SELECT * FROM reservations"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*entity.Reservation
	for rows.Next() {
		var reservation entity.Reservation
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *reservationRepo) GetByUserID(ctx context.Context, userID int) ([]*entity.Reservation, error) {
	query := "SELECT * FROM reservations WHERE user_id = $1"
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*entity.Reservation
	for rows.Next() {
		var reservation entity.Reservation
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}
