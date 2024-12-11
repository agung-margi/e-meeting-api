package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/presenter/model"
)

type userUseCase struct {
	repo repository.UserRepository
}

func NewStudentUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *userUseCase) UpsertUser(ctx context.Context, user *model.User) error {
	dataUser := entity.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		ImgUrl:   user.ImgUrl,
	}
	return u.repo.SaveUser(ctx, &dataUser)
}
