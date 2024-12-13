package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/presenter/model"
	"errors"
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

func (u *userUseCase) UpsertUser(ctx context.Context, user *model.UserRequest) error {
	emailExists, err := u.repo.CheckEmailExists(ctx, user.Email)
	if err != nil {
		return err
	}

	if emailExists {
		return errors.New("email already exists")
	}

	dataUser := entity.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	return u.repo.SaveUser(ctx, &dataUser)
}
