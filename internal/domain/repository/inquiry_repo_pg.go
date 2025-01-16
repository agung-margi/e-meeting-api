package repository

import (
	"context"
	"database/sql"
	"e-meeting-api/internal/domain/entity"
)

type inquiryRepo struct {
	DB *sql.DB
}

func NewInquiryRepository(db *sql.DB) InquiryRepository {
	return &inquiryRepo{db}
}

func (r *inquiryRepo) Save(ctx context.Context, inquiry *entity.Inquiry) (*entity.Inquiry, error) {
	query := `
		INSERT INTO inquiries (
			user_id, room_id, room_name, room_type, room_capacity, room_price,
			date_reservation, start_time, end_time, duration, name, phone, company,
			participants, snack_id, snack_name, snack_category, snack_price,
			total_room_price, total_snack_price, total_price, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24
		)
		RETURNING	id
		`
	var id int
	err := r.DB.QueryRowContext(ctx, query,
		inquiry.UserID, inquiry.RoomID, inquiry.RoomName, inquiry.RoomType, inquiry.RoomCapacity, inquiry.RoomPrice,
		inquiry.BookingDate, inquiry.StartTime, inquiry.EndTime, inquiry.Duration, inquiry.Name, inquiry.Phone, inquiry.Company,
		inquiry.Participants, inquiry.SnackID, inquiry.SnackName, inquiry.SnackCategory, inquiry.SnackPrice,
		inquiry.TotalRoomPrice, inquiry.TotalSnackPrice, inquiry.TotalPrice, inquiry.Notes, inquiry.CreatedAt, inquiry.UpdatedAt).Scan(&id)

	if err != nil {
		return nil, err
	}
	inquiry.ID = id
	return inquiry, nil
}

func (r *inquiryRepo) GetByID(ctx context.Context, id int, userId int) (*entity.Inquiry, error) {
	query := "SELECT * FROM inquiries WHERE id = $1 AND user_id = $2"
	inquiry := &entity.Inquiry{}
	err := r.DB.QueryRowContext(ctx, query, id, userId).Scan(&inquiry.ID, &inquiry.UserID, &inquiry.RoomID, &inquiry.RoomName, &inquiry.RoomType, &inquiry.RoomCapacity, &inquiry.RoomPrice,
		&inquiry.BookingDate, &inquiry.StartTime, &inquiry.EndTime, &inquiry.Duration, &inquiry.Name, &inquiry.Phone, &inquiry.Company,
		&inquiry.Participants, &inquiry.SnackID, &inquiry.SnackName, &inquiry.SnackCategory, &inquiry.SnackPrice,
		&inquiry.TotalRoomPrice, &inquiry.TotalSnackPrice, &inquiry.TotalPrice, &inquiry.Notes, &inquiry.CreatedAt, &inquiry.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return inquiry, nil
}

func (r *inquiryRepo) DeleteInquiry(ctx context.Context, id int, userId int) error {
	query := "DELETE FROM inquiries WHERE id = $1 AND user_id = $2"
	_, err := r.DB.ExecContext(ctx, query, id, userId)
	if err != nil {
		return err
	}
	return nil
}
