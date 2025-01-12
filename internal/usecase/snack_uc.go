package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
)

type snackUseCase struct {
	repo repository.SnackRepository
}

func NewSnackUseCase(repo repository.SnackRepository) SnackUseCase {
	return &snackUseCase{repo: repo}
}

func (u *snackUseCase) GetAll(ctx context.Context) ([]entity.Snack, error) {
	return u.repo.GetAll(ctx)
}

func (u *snackUseCase) GetByID(ctx context.Context, id int) (*entity.Snack, error) {
	return u.repo.GetByID(ctx, id)
}
