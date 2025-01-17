package usecase

import (
	"context"
	"e-meeting-api/configs"
	"e-meeting-api/internal/domain/entity"
	"e-meeting-api/internal/domain/repository"
	"e-meeting-api/pkg/util"
	"e-meeting-api/presenter/model"
	"errors"
	"fmt"
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

	userExists, err := u.repo.CheckUsernameExists(ctx, user.Username)
	if err != nil {
		return err
	}
	if userExists {
		return errors.New("username already exists")
	}

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
	// Ambil data user lama dari database
	existingUser, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Periksa properti dari request dan gunakan data lama jika properti kosong
	updatedUser := &entity.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
		ImgUrl:   user.ImgUrl,
		IsActive: user.IsActive,
		Language: user.Language,
	}

	if updatedUser.Username == "" {
		updatedUser.Username = existingUser.Username
	}
	if updatedUser.Email == "" {
		updatedUser.Email = existingUser.Email
	}
	if updatedUser.Password == "" {
		updatedUser.Password = existingUser.Password
	}
	if updatedUser.ImgUrl == "" {
		updatedUser.ImgUrl = existingUser.ImgUrl
	}
	if updatedUser.Language == "" {
		updatedUser.Language = existingUser.Language
	}
	// Jika IsActive tidak di-set dalam request, gunakan nilai lama
	if user.IsActive == false && !util.IsFieldProvided(user, "IsActive") {
		updatedUser.IsActive = existingUser.IsActive
	}
	// Jika IsAdmin tidak di-set dalam request, gunakan nilai lama
	if user.IsAdmin == false && !util.IsFieldProvided(user, "IsAdmin") {
		updatedUser.IsAdmin = existingUser.IsAdmin
	}

	// Kirim data yang sudah diperbarui ke repository
	return u.repo.UpdateUser(ctx, id, updatedUser)
}

func (u *userUseCase) Login(ctx context.Context, username string, password string) (string, *entity.User, error) {

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
		"exp":      time.Now().Add(6 * time.Hour).Unix(),
	}

	token, err := util.GenerateToken(6*time.Hour, claims)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (u *userUseCase) GetResetPassword(ctx context.Context, email string) error {
	//getname
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("email not found")
	}

	name := user.Username

	token := util.GenerateRandomString(20)
	expiredAt := time.Now().Add(30 * time.Minute)

	err = u.repo.SaveResetToken(ctx, email, token, expiredAt)
	if err != nil {
		return err
	}

	baseURLFE := configs.AppConfig.BaseURLFE
	resetLink := fmt.Sprintf("%s/reset-password/%s?email=%s", baseURLFE, token, email)

	err = util.SendResetPasswordEmail(email, name, resetLink)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) ResetPassword(ctx context.Context, email string, token string, newPassword string) error {
	// Validate the reset token
	isValid, err := u.repo.ValidateResetToken(ctx, email, token)
	if err != nil || !isValid {
		return errors.New("Invalid or expired reset token")
	}

	// Fetch the user by email
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("Error fetching user data")
	}
	if user == nil {
		return errors.New("User not found")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed to hash password")
	}

	// Update the password in the database
	err = u.repo.UpdatePassword(ctx, user.ID, string(hashedPassword))
	if err != nil {
		return errors.New("Failed to update password")
	}

	// Delete the reset token
	err = u.repo.DeleteResetToken(ctx, email)
	if err != nil {
		return errors.New("Failed to delete reset token")
	}

	return nil
}
