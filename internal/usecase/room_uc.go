package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/presenter/model"
	"log"
)

type roomUseCase struct {
	repo repository.RoomRepository
}

func NewRoomUseCase(repo repository.RoomRepository) RoomUseCase {
	return &roomUseCase{repo: repo}
}

func (u *roomUseCase) GetAll(ctx context.Context) ([]entity.Room, error) {
	return u.repo.GetAll(ctx)
}

func (u *roomUseCase) GetByID(ctx context.Context, id int) (*entity.RoomWithType, error) {
	return u.repo.GetByID(ctx, id) // Return the correct type
}

func (u *roomUseCase) SaveRoom(ctx context.Context, roomRequest *model.RoomRequest) error {
	room := &entity.Room{
		Name:     roomRequest.Name,
		RoomType: roomRequest.RoomType,
		Price:    roomRequest.Price,
		Capacity: roomRequest.Capacity,
		ImgUrl:   roomRequest.ImgUrl,
	}

	err := u.repo.SaveRoom(ctx, room)
	if err != nil {
		log.Printf("Error saving room: %v", err)
		return err
	}
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
		log.Printf("Error updating room: %v", err)
		return err
	}
	return nil
}
