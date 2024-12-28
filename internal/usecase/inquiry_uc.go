package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"errors"
	"fmt"
	"strings"
	"time"
)

type inquiryUsecase struct {
	inquiryRepo repository.InquiryRepository
	roomRepo    repository.RoomRepository
	snackRepo   repository.SnackRepository
}

func NewInquiryUseCase(inquiryRepo repository.InquiryRepository, roomRepo repository.RoomRepository, snackRepo repository.SnackRepository) InquiryUsecase {
	return &inquiryUsecase{inquiryRepo: inquiryRepo, roomRepo: roomRepo, snackRepo: snackRepo}
}

func (u *inquiryUsecase) Save(ctx context.Context, reservation *entity.Reservation) (*entity.Inquiry, error) {
	//get room details
	room, err := u.roomRepo.GetByID(ctx, reservation.RoomID)
	if err != nil {
		return nil, err
	}

	//get snack details
	var snack *entity.Snack
	if reservation.SnackID != nil {
		fmt.Println("Fetching snack with ID:", *reservation.SnackID)
		snack, err = u.snackRepo.GetByID(ctx, *reservation.SnackID)
		if err != nil {
			return nil, err
		}
		fmt.Println("Fetched snack:", snack)
	}
	date := reservation.BookingDate
	start := reservation.StartTime
	end := reservation.EndTime
	format := "2006-01-02 15:00"

	// format := "2006-01-02 15:00"
	// date := req.BookingDate
	startStr := fmt.Sprintf("%s %s", date, start)
	endStr := fmt.Sprintf("%s %s", date, end)

	startTime, err := time.Parse(format, startStr)
	if err != nil {
		return nil, err
	}

	// Parse end time
	endTime, err := time.Parse(format, endStr)
	if err != nil {
		return nil, err
	}

	if startTime.After(endTime) {
		return nil, errors.New("i start time should be before end time")
	}

	if startTime.Before(time.Now()) {
		return nil, errors.New("start time should be in the future")
	}

	if reservation.Participants < 1 {
		return nil, errors.New("participants should be at least 1")

	}
	if reservation.Participants > room.Capacity {
		return nil, errors.New("room capacity is not enough")
	}

	// Calculate duration
	duration := int(endTime.Sub(startTime).Hours())

	roomTotalPrice := room.Price * duration
	snackTotalPrice := 0
	grandTotal := roomTotalPrice

	//create inquiry
	inquiry := &entity.Inquiry{
		UserID:          reservation.UserID,
		RoomID:          reservation.RoomID,
		RoomName:        room.Name,
		RoomType:        room.RoomType,
		RoomCapacity:    room.Capacity,
		RoomPrice:       room.Price,
		BookingDate:     reservation.BookingDate,
		StartTime:       startStr,
		EndTime:         endStr,
		Duration:        duration,
		Name:            reservation.Name,
		Phone:           reservation.Phone,
		Company:         reservation.Company,
		Participants:    reservation.Participants,
		SnackID:         nil,
		SnackName:       "",
		SnackCategory:   "",
		SnackPrice:      0,
		TotalRoomPrice:  roomTotalPrice,
		TotalSnackPrice: snackTotalPrice,
		TotalPrice:      grandTotal,
		Notes:           reservation.Notes,
		CreatedAt:       strings.Split(time.Now().String(), ".")[0],
		UpdatedAt:       strings.Split(time.Now().String(), ".")[0],
	}

	fmt.Println(inquiry)
	// If snack details were found, update the Inquiry struct
	if snack != nil {
		inquiry.SnackPrice = snack.Price
		inquiry.SnackID = reservation.SnackID
		inquiry.SnackName = snack.Name
		inquiry.SnackCategory = snack.Category
		inquiry.SnackPrice = snack.Price
		inquiry.TotalSnackPrice = snack.Price * reservation.Participants
		inquiry.TotalPrice = roomTotalPrice + snackTotalPrice
	}
	inquiry, err = u.inquiryRepo.Save(ctx, inquiry)
	if err != nil {
		fmt.Println("Error saving inquiry:", err)
		return nil, err
	}

	inquiry.ID = inquiry.ID
	return inquiry, nil
}
