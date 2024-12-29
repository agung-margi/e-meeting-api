package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/presenter/model"
	"errors"
	"time"
)

type reservationRepo struct {
	DB *sql.DB
}

func NewReservationRepository(db *sql.DB) ReservationRepository {
	return &reservationRepo{DB: db}
}

func (r *reservationRepo) GetByID(ctx context.Context, id int) (*entity.Reservation, error) {
	reservation := &entity.Reservation{}
	query := `
    SELECT id, room_id,user_id, start_time, end_time, status, created_at, updated_at
    FROM reservations
    WHERE id = $1
`
	row := r.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&reservation.ID, &reservation.RoomID, &reservation.UserID, &reservation.StartTime, &reservation.EndTime, &reservation.Status, &reservation.CreatedAt, &reservation.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}
	return reservation, nil
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
func (r *reservationRepo) GetReservationsByRoomAndDate(ctx context.Context, roomID int, date string) ([]entity.RoomSchedule, error) {
	query := `
			SELECT r.id, r.room_id, rm.name, r.start_time, r.end_time 
FROM reservations r
JOIN rooms rm ON r.room_id = rm.id
WHERE r.room_id = $1
  AND DATE(r.start_time) = $2
  AND (r.status = 'booked' OR r.status = 'paid')
	`

	rows, err := r.DB.QueryContext(ctx, query, roomID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations := make([]entity.RoomSchedule, 0)
	for rows.Next() {
		reservation := &entity.RoomSchedule{}
		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.RoomName, &reservation.StartTime, &reservation.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, *reservation)
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
func (r *reservationRepo) Save(ctx context.Context, reservation *entity.Reservation, details *entity.ReservationDetails) (*entity.Reservation, error) {
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
		return nil, err
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
		return nil, err
	}

	return reservation, nil
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

func (r *reservationRepo) GetAll(ctx context.Context, startDate, endDate *time.Time, roomType int, status string, userID *int) ([]*entity.Reservation, error) {
	query := `
		SELECT
    r.id,
    r.user_id,
    r.room_id,
    r.start_time,
    r.end_time,
    r.booking_date,
    r.room_price,
    r.snack_price,
    r.total_price,
    r.status,
    rt.name AS room_type,
    r.created_at,
    r.updated_at
FROM
    reservations r
JOIN
    rooms rm
ON
    r.room_id = rm.id
JOIN
    room_types rt
ON
    rm.room_type_id = rt.id
WHERE
    ($1::timestamp IS NULL OR r.start_time >= $1::timestamp)
    AND ($2::timestamp IS NULL OR r.end_time <= $2::timestamp)
    AND ($3::int = 0 OR rm.room_type_id = $3::int)
    AND ($4::text = '' OR r.status = $4::text)
		AND ($5::int IS NULL OR r.user_id = $5::int)
		`
	// Handle parameter nil
	var params []interface{}
	params = append(params, startDate, endDate, roomType, status, userID)
	rows, err := r.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parsing hasil query ke dalam slice of Reservation
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
			&reservation.RoomType,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	return reservations, nil
}

func (r *reservationRepo) GetRoomSchedule(ctx context.Context, roomID int, date string) error {
	query := " "
	_, err := r.DB.ExecContext(ctx, query, roomID, date)
	return err
}

func (r *reservationRepo) GetReservationDetails(ctx context.Context, reservationID int) (*[]model.ReservationDetailsResponse, error) {
	query := `
		SELECT 
			r.user_id AS user_id,
			rm.id AS room_id,
			rm.name AS room_name,
			rm.capacity AS room_capacity,
			rm.price AS room_price,

			-- Booking Details
			rd.name AS customer_name,
			rd.phone AS customer_phone,
			rd.company AS company_name,
			r.booking_date AS date_reservation,
			r.start_time AS start_date,
			r.end_time AS end_date,
			rd.participants AS participants,

			-- Consumption Details
			COALESCE(s.id, 0) AS snack_id,
			COALESCE(s.category, '') AS snack_category,
			COALESCE(s.name, '') AS snack_name,
			COALESCE(s.price, 0) AS snack_price,

			-- Total and Notes
			r.total_price AS total_price,
			r.room_price AS total_room_price,
			COALESCE(r.snack_price, 0) AS total_snack_price,
			rd.notes AS reservation_notes

		FROM 
			reservations r
		JOIN 
			rooms rm ON r.room_id = rm.id
		JOIN 
			reservation_details rd ON r.id = rd.reservation_id
		LEFT JOIN 
			snacks s ON rd.snack_id = s.id
		WHERE 
			r.id = $1;
	`

	row := r.DB.QueryRowContext(ctx, query, reservationID)

	// Create response structure
	reservation := &model.ReservationDetailsResponse{}

	// Variables to hold query results
	var (
		roomID, roomCapacity, roomPrice, participants, snackID, snackPrice, totalPrice, totalRoomPrice, totalSnackPrice int
		roomName, customerName, customerPhone, companyName, snackCategory, snackName, notes                             string
		userID                                                                                                          int
		dateReservation, startDate, endDate                                                                             time.Time
	)

	// Scan the results from the query
	err := row.Scan(
		&userID,
		&roomID, &roomName, &roomCapacity, &roomPrice,
		&customerName, &customerPhone, &companyName, &dateReservation, &startDate, &endDate, &participants,
		&snackID, &snackCategory, &snackName, &snackPrice,
		&totalPrice, &totalRoomPrice, &totalSnackPrice,
		&notes,
	)
	if err != nil {
		return nil, err
	}

	// Calculate duration
	duration := int(endDate.Sub(startDate).Hours())

	// Populate RoomDetails
	reservation.RoomDetails = []model.RoomDetailsResponse{
		{
			RoomID:    roomID,
			RoomName:  roomName,
			Capacity:  roomCapacity,
			RoomPrice: roomPrice,
		},
	}

	// Populate BookDetails
	reservation.BookDetails = []model.BookingDetailsResponse{
		{
			Name:            customerName,
			Phone:           customerPhone,
			Company:         companyName,
			DateReservation: dateReservation.Format("2006-01-02"),
			Duration:        duration,
			Participants:    participants,
		},
	}

	// Populate Consumption
	if snackID != 0 {
		reservation.SnacksDetails = []model.SnackDetailsResponse{
			{
				SnackID:  &snackID,
				Category: snackCategory,
				Name:     snackName,
				Price:    snackPrice,
			},
		}
	}

	// Populate TotalPrice
	reservation.TotalPrice = []model.ReservationTotalPriceResponse{
		{
			RoomName:        roomName,
			Duration:        duration,
			RoomPrice:       roomPrice,
			TotalRoomPrice:  totalRoomPrice,
			SnackName:       snackName,
			Qty:             participants,
			TotalSnackPrice: snackPrice * participants,
			TotalPrice:      snackPrice*participants + roomPrice*duration,
		},
	}

	// Add Notes
	reservation.Notes = notes
	reservation.ReservationID = reservationID
	reservation.UserID = userID

	return &[]model.ReservationDetailsResponse{*reservation}, nil
}

func (r *reservationRepo) UpdateStatus(ctx context.Context, reservationID int, status string) error {
	query := "UPDATE reservations SET status = $1 WHERE id = $2"
	_, err := r.DB.ExecContext(ctx, query, status, reservationID)
	return err
}
