package usecase

import (
	"context"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/pkg/util"
	"e-meeting-api/presenter/model"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	repo repository.UserRepository
}

func NewStudentUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) GetUser(ctx context.Context, id int) (*entity.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *userUseCase) SaveUser(ctx context.Context, user *model.UserRequest) error {

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

func (u *userUseCase) UpdateUser(ctx context.Context, id int, user *model.UpdateUserRequest) error {
	existingUser, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	updatedUser := &entity.User{
		Password:  user.Password,
		ImgUrl:    user.ImgUrl,
		IsActive:  user.IsActive,
		UpdatedAt: time.Now(),
	}
	return u.repo.UpdateUser(ctx, id, updatedUser)
}

func (u *userUseCase) Login(ctx context.Context, username string, password string) (string, *entity.User, error) {
	secret := os.Getenv("JWT_SECRET")

	//Check username exists
	user, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		return "", nil, errors.New("username or password invalid")
	}

	//Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("username or password invalid")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"language": user.Language,
		"exp":      time.Now().Add(12 * time.Hour).Unix(),
	}

	token, err := util.GenerateToken(secret, 24*time.Hour, claims)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
