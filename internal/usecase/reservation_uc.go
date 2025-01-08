package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/presenter/model"
	"errors"
	"fmt"
	"time"
)

type reservationUseCase struct {
	inquiryRepo repository.InquiryRepository
	repo        repository.ReservationRepository
}

func NewReservationUseCase(repo repository.ReservationRepository, inquiryRepo repository.InquiryRepository) ReservationUseCase {
	return &reservationUseCase{
		repo:        repo,
		inquiryRepo: inquiryRepo}
}

func (u *reservationUseCase) GetReservation(ctx context.Context, reservationID int) (*[]model.ReservationDetailsResponse, error) {
	return u.repo.GetReservationDetails(ctx, reservationID)
}

func (u *reservationUseCase) GetReservationsByRoomAndDate(ctx context.Context, roomId int, date string) ([]entity.RoomSchedule, error) {
	return u.repo.GetReservationsByRoomAndDate(ctx, roomId, date)
}
func (u *reservationUseCase) PayReservation(ctx context.Context, id int, userId int, isAdmin bool) error {

	reservation, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if reservation.UserID != userId && !isAdmin {
		return errors.New("invalid access")
	}

	// validasi status pembayaran
	if reservation.Status == "paid" {
		return errors.New("reservation that has already been paid")
	}

	// if time.Now().After(reservation.StartTime) {
	// 	return errors.New("reservation that has already started")
	// }

	// if time.Now().After(reservation.EndTime) {
	// 	return errors.New("reservation that has already ended")
	// }

	if reservation.Status == "cancelled" {
		return errors.New("reservation that has already been cancelled")
	}

	return u.repo.UpdateStatus(ctx, id, "paid")
}

func (u *reservationUseCase) CancelReservation(ctx context.Context, id int, userId int, isAdmin bool) error {
	reservation, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if reservation.UserID != userId && !isAdmin {
		return errors.New("unauthorized")
	}

	if reservation.Status == "paid" {
		return errors.New("reservation that has already been paid")
	}

	// if time.Now().After(reservation.StartTime) {
	// 	return errors.New("reservation that has already started")
	// }

	// if time.Now().After(reservation.EndTime) {
	// 	return errors.New("reservation that has already ended")
	// }

	if reservation.Status == "cancelled" {
		return errors.New("reservation that has already been cancelled")
	}

	return u.repo.UpdateStatus(ctx, id, "cancelled")
}

func (u *reservationUseCase) Save(ctx context.Context, inquiryId int, userId int) (*[]model.ReservationDetailsResponse, error) {
	if u.inquiryRepo == nil {
		return nil, errors.New("inquiry repository is not initialized")
	}

	inquiry, err := u.inquiryRepo.GetByID(ctx, inquiryId, userId)

	if err != nil {
		return nil, errors.New("inquiry not found")
	}

	if inquiry == nil {
		return nil, errors.New("inquiry not found")
	}

	// Cek ketersediaan ruangan
	available, err := u.repo.CheckAvailability(ctx, inquiry.RoomID, inquiry.StartTime, inquiry.EndTime)
	if err != nil {
		return nil, err
	}

	// Jika ruangan tidak tersedia, return error
	if !available {
		return nil, errors.New("room is not available for the selected time")
	}

	// Buat reservation berdasarkan inquiry
	reservation := &entity.Reservation{
		UserID:      userId,
		RoomID:      inquiry.RoomID,
		StartTime:   inquiry.StartTime,
		EndTime:     inquiry.EndTime,
		BookingDate: inquiry.BookingDate,
		RoomPrice:   inquiry.RoomPrice,
		TotalPrice:  inquiry.TotalPrice,
	}

	// Panggil fungsi Save pada repository
	reservation, err = u.repo.Save(ctx, reservation, &entity.ReservationDetails{
		ReservationID: reservation.ID,
		Name:          inquiry.Name,
		Phone:         inquiry.Phone,
		Company:       inquiry.Company,
		SnackID:       inquiry.SnackID,
		Participants:  inquiry.Participants,
		Notes:         inquiry.Notes,
	})
	if err != nil {
		return nil, err
	}
	reservationDetails, err := u.repo.GetReservationDetails(ctx, reservation.ID)
	fmt.Println("reservationDetails:", reservationDetails)
	if err != nil {
		return nil, err
	}
	return reservationDetails, nil
}
func (u *reservationUseCase) GetAll(ctx context.Context, startDate *time.Time, endDate *time.Time, roomType int, status string, userID *int) ([]*entity.ReservationHistory, error) {
	reservations, err := u.repo.GetAll(ctx, startDate, endDate, roomType, status, userID)
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (u *reservationUseCase) GetSchedulesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.RoomSchedule, error) {
	return u.repo.GetSchedulesByDateRange(ctx, startDate, endDate)
}

func (u *reservationUseCase) GetDashboardData(ctx context.Context, startDate, endDate time.Time) (entity.Dashboard, error) {
	// Ambil data dashboard dan room omzet dari repository
	dashboard, err := u.repo.GetDashboardDataByDateRange(ctx, startDate, endDate)
	if err != nil {
		return entity.Dashboard{}, err
	}

	roomOmzet, err := u.repo.GetRoomOmzetByDateRange(ctx, startDate, endDate)
	if err != nil {
		return entity.Dashboard{}, err
	}

	// Hitung presentase omzet per room
	for i := range roomOmzet {
		if dashboard.TotalOmzet > 0 {
			roomOmzet[i].Presentage = (float64(roomOmzet[i].RoomOmzet) / float64(dashboard.TotalOmzet)) * 100
		} else {
			roomOmzet[i].Presentage = 0
		}
	}

	// Gabungkan dashboard dan roomOmzet ke dalam satu respons
	response := entity.Dashboard{
		TotalReservation: dashboard.TotalReservation,
		TotalSnackPrice:  dashboard.TotalSnackPrice,
		TotalPrice:       dashboard.TotalPrice,
		TotalRoom:        dashboard.TotalRoom,
		TotalVisitor:     dashboard.TotalVisitor,
		TotalOmzet:       dashboard.TotalOmzet,
		RoomDetails:      roomOmzet,
	}

	return response, nil
}
