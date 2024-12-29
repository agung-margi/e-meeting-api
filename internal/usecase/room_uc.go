package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/presenter/model"
	"fmt"
)

type roomUseCase struct {
	repo repository.RoomRepository
}

func NewRoomUseCase(repo repository.RoomRepository) RoomUseCase {
	return &roomUseCase{repo: repo}
}

func (u *roomUseCase) GetAll(ctx context.Context, name string, roomType int, capacity int) ([]entity.Room, error) {
	rooms, err := u.repo.GetAll(ctx, name, roomType, capacity)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (u *roomUseCase) GetByID(ctx context.Context, id int) (*entity.RoomWithType, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *roomUseCase) SaveRoom(ctx context.Context, roomRequest *model.RoomRequest) error {

	// Creating the Room entity from RoomRequest.
	room := &entity.Room{
		Name:     roomRequest.Name,
		RoomType: roomRequest.RoomType,
		Price:    roomRequest.Price,
		Capacity: roomRequest.Capacity,
		ImgUrl:   roomRequest.ImgUrl,
	}

	err := u.repo.SaveRoom(ctx, room)
	if err != nil {
		return err
	}

	// Debugging: Log a success message when the room is successfully saved.
	fmt.Println("Room saved successfully.")

	return nil
}

func (u *roomUseCase) UpdateRoom(ctx context.Context, id int, roomRequest *model.RoomRequest) error {

	room := &entity.Room{
		ID:       id,
		Name:     roomRequest.Name,
		RoomType: roomRequest.RoomType,
		Price:    roomRequest.Price,
		Capacity: roomRequest.Capacity,
		ImgUrl:   roomRequest.ImgUrl,
	}

	err := u.repo.UpdateRoom(ctx, id, room)
	if err != nil {
		return err
	}
	return nil
}

func (u *roomUseCase) DeleteRoom(ctx context.Context, id int) error {
	return u.repo.DeleteRoom(ctx, id)
}
