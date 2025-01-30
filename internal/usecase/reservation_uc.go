package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/presenter/model"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/xuri/excelize/v2"
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

	startTime, err := time.Parse(time.RFC3339, reservation.StartTime)
	if err != nil {
		return errors.New("invalid start time format")
	}

	endTime, err := time.Parse(time.RFC3339, reservation.EndTime)
	if err != nil {
		return errors.New("invalid end time format")
	}
	if time.Now().After(startTime) {
		return errors.New("reservation that has already started")
	}

	if time.Now().After(endTime) {
		return errors.New("reservation that has already ended")
	}

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

	startTime, err := time.Parse(time.RFC3339, reservation.StartTime)
	if err != nil {
		return errors.New("invalid start time format")
	}

	endTime, err := time.Parse(time.RFC3339, reservation.EndTime)
	if err != nil {
		return errors.New("invalid end time format")
	}
	if time.Now().After(startTime) {
		return errors.New("reservation that has already started")
	}

	if time.Now().After(endTime) {
		return errors.New("reservation that has already ended")
	}

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
		RoomPrice:   inquiry.TotalRoomPrice,
		SnackPrice:  inquiry.TotalSnackPrice,
		TotalPrice:  inquiry.TotalPrice,
		Status:      "booked",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(24 * time.Hour),
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

	//Delete inquiry setelah berhasil
	if err := u.inquiryRepo.DeleteInquiry(ctx, inquiryId, userId); err != nil {
		return nil, err
	}
	reservationDetails, err := u.repo.GetReservationDetails(ctx, reservation.ID)
	if err != nil {
		return nil, err
	}
	return reservationDetails, nil
}
func (u *reservationUseCase) GetAll(ctx context.Context, startDate string, endDate string, roomType int, status string, userID *int) ([]*entity.ReservationHistory, error) {
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

func (u *reservationUseCase) ExpireReservations(ctx context.Context) error {
	expiredReservations, err := u.repo.GetExpiredReservations(ctx, time.Now())
	if err != nil {
		return err
	}

	for _, reservation := range expiredReservations {
		err := u.repo.UpdateStatus(ctx, reservation.ID, "cancelled")
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *reservationUseCase) StartExpiredReservationWorker(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := u.ExpireReservations(ctx)
			if err != nil {
				log.Printf("Error expiring reservations: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (u *reservationUseCase) ExportReservations(ctx context.Context, startDate string, endDate string, roomType int, status string, userID *int) (*excelize.File, error) {

	// if startDate == nil || endDate == nil {
	// 	return nil, fmt.Errorf("start_date and end_date are required")
	// }
	// if startDate.After(*endDate) {
	// 	return nil, fmt.Errorf("start_date must be before end_date")
	// }

	// Validasi status yang diperbolehkan
	if status != "" && status != "paid" && status != "booked" && status != "cancel" {
		return nil, fmt.Errorf("invalid status")
	}

	// Ambil data reservasi dari database
	reservations, err := u.repo.GetAll(ctx, startDate, endDate, roomType, status, userID)
	if err != nil {
		return nil, err
	}

	fmt.Println("reservations: %v\n", startDate, endDate, roomType, status, userID)

	if reservations == nil {
		return nil, fmt.Errorf("reservations not found")
	}

	// Buat file Excel

	f := excelize.NewFile()
	sheetName := "Reservations"
	f.SetSheetName("Sheet1", sheetName)

	// Header kolom
	headers := []string{"Booking Date", "Room Name", "Room Type", "Status"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Data reservasi
	for i, reservation := range reservations {
		row := i + 2
		bookingDate, err := time.Parse("2006-01-02T15:04:05Z", reservation.BookingDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse booking date: %v", err)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), bookingDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), reservation.RoomName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), reservation.RoomType)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), reservation.Status)
	}

	return f, nil
}
