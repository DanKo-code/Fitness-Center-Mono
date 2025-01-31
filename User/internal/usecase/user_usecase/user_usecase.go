package user_usecase

import (
	"User/internal/dtos"
	customErrors "User/internal/errors"
	"User/internal/models"
	"User/internal/repository"
	"User/internal/usecase"
	"User/pkg/logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserUseCase struct {
	userRepo     repository.UserRepository
	cloudUseCase usecase.CloudUseCase
}

func NewUserUseCase(
	userRepo repository.UserRepository,
	cloudUseCase usecase.CloudUseCase,
) *UserUseCase {
	return &UserUseCase{
		userRepo:     userRepo,
		cloudUseCase: cloudUseCase,
	}
}

func (u *UserUseCase) CreateUser(
	ctx context.Context,
	cmd *dtos.CreateUserCommand,
) (*models.User, error) {

	hashedPassword, err := HashPassword(cmd.Password)
	if err != nil {
		logger.ErrorLogger.Printf("Error hashing password: %v", err)
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Name:         cmd.Name,
		Email:        cmd.Email,
		Role:         cmd.Role,
		PasswordHash: hashedPassword,
		Photo:        "",
		UpdatedTime:  time.Now(),
		CreatedTime:  time.Now(),
	}

	createUser, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return createUser, nil
}

func (u *UserUseCase) GetUserById(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) UpdateUser(
	ctx context.Context,
	cmd *dtos.UpdateUserCommand,
) (*models.User, error) {

	err := u.userRepo.UpdateUser(ctx, cmd)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetUserById(ctx, cmd.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) DeleteUserById(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {

	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return nil, customErrors.UserNotFound
	}

	err = u.userRepo.DeleteUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.Photo != "" {
		prefix := "user/"
		index := strings.Index(user.Photo, prefix)
		var s3PhotoKey string
		if index != -1 {
			s3PhotoKey = user.Photo[index+len(prefix):]
		} else {
			logger.ErrorLogger.Printf("Prefix not found")
			return nil, fmt.Errorf("prefix not found")
		}

		err = u.cloudUseCase.DeleteObject(ctx, "user/"+s3PhotoKey)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (u *UserUseCase) GetUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) CheckPassword(
	ctx context.Context,
	cmd *dtos.CheckPasswordCommand,
) error {

	user, err := u.userRepo.GetUserById(ctx, cmd.UserId)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cmd.Password)); err != nil {
		logger.ErrorLogger.Printf("Error Invalid Password: %v", err)
		return customErrors.InvalidPassword
	}

	return nil
}

func (u *UserUseCase) GetUsersByIds(ctx context.Context, ids []uuid.UUID) ([]*models.User, error) {
	resp, err := u.userRepo.GetUsersByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (u *UserUseCase) GetClients(ctx context.Context) ([]*models.User, error) {
	clients, err := u.userRepo.GetClients(ctx)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (u *UserUseCase) GetAdmins(ctx context.Context) ([]*models.User, error) {
	admins, err := u.userRepo.GetAdmins(ctx)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
