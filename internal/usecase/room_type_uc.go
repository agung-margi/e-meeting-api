package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
)

type roomtypeUseCase struct {
	repo repository.RoomtypeRepository
}

func NewRoomtypeUseCase(repo repository.RoomtypeRepository) RoomtypeUseCase {
	return &roomtypeUseCase{repo: repo}
}

func (u *roomtypeUseCase) GetAll(ctx context.Context) ([]entity.RoomtypeRepository, error) {
	return u.repo.GetAll(ctx)
}

func (u *roomtypeUseCase) GetByID(ctx context.Context, id int) (*entity.RoomtypeRepository, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *roomtypeUseCase) GetByName(ctx context.Context, name string) (*entity.RoomtypeRepository, error) {
	return u.repo.GetByName(ctx, name)
}

func (u *roomtypeUseCase) SaveRoomType(ctx context.Context, roomType *entity.RoomtypeRepository) error {
	return u.repo.SaveRoomType(ctx, roomType)
}
