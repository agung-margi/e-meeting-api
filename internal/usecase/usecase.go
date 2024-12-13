package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/presenter/model"
)

type UserUseCase interface {
	GetUser(ctx context.Context, id string) (*entity.User, error)
	UpsertUser(ctx context.Context, user *model.UserRequest) error
}
