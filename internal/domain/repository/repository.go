package repository

import (
	"context"
	"e-meeting-api/internal/domain/entity"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	SaveUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	Login(ctx context.Context, username string, password string) (*entity.User, error)
}
