package service

import (
	"context"
	"database/sql"
	"test/api/entity"
	"test/api/model"
	"test/api/repository"
	"test/api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	RegisterUser(ctx context.Context, user model.UserRegisterResponse, signature string, expired int) (model.UserRegisterResponse, error)
}

type UserServiceImplement struct {
	userRepository    repository.UserRepository
	emailVerification repository.EmailVerificationRepository
	DB                *sql.DB
}

func NewUserService(userRepository *repository.UserRepository, db *sql.DB, emailVerification *repository.EmailVerificationRepository) UserServiceImplement {
	return UserServiceImplement{
		userRepository:    *userRepository,
		emailVerification: *emailVerification,
		DB:                db,
	}
}

// RegisterUser is used to register new user
func (service *UserServiceImplement) RegisterUser(ctx *gin.Context, user model.UserRegisterResponse, signature string, expired int) (model.UserRegisterResponse, error) {
	var response model.UserRegisterResponse

	tx, err := service.DB.Begin()

	if err != nil {
		return model.UserRegisterResponse{}, err
	}
	defer utils.CommitOrRollback(tx)

	user.Created_at = time.Now()
	user.Updated_at = time.Now()
	user.EmailVerification = time.Now()

	err = service.userRepository.Register(ctx, tx, entity.Users{
		Name:              user.Name,
		Username:          user.Username,
		Email:             user.Email,
		Password:          user.Password,
		Role:              user.Role,
		Phone:             user.Phone,
		Gender:            user.Gender,
		DisabilityType:    user.DisabilityType,
		Birthdate:         user.Birthdate,
		EmailVerification: user.EmailVerification,
		CreatedAt:         user.Created_at,
		UpdatedAt:         user.Updated_at,
	})
	if err != nil {
		return model.UserRegisterResponse{}, err
	}
	temp, err := service.userRepository.GetLastInsertUser(ctx, tx)
	if err != nil {
		return model.UserRegisterResponse{}, err
	}

	response = model.UserRegisterResponse{
		Id:                temp.Id,
		Name:              temp.Name,
		Username:          temp.Username,
		Email:             temp.Email,
		Password:          temp.Password,
		Role:              temp.Role,
		Phone:             temp.Phone,
		Gender:            temp.Gender,
		DisabilityType:    user.DisabilityType,
		Birthdate:         user.Birthdate,
		EmailVerification: temp.EmailVerification,
		Created_at:        temp.CreatedAt,
		Updated_at:        temp.UpdatedAt,
	}

	// Send Email Verification
	// Check If email is exists in table email verifications
	rows, err := service.emailVerification.FindByEmail(ctx, tx, user.Email)
	if err != nil {
		return model.UserRegisterResponse{}, err
	}
	// Create new signature if not exist
	emailVerification := entity.EmailVerification{
		Email:     user.Email,
		Signature: signature,
		Expired:   expired,
	}
	if rows.Email == "" {
		_, err = service.emailVerification.Create(ctx, tx, emailVerification)
		if err != nil {
			return model.UserRegisterResponse{}, err
		}
	} else {
		// Update token if email is exist
		_, err = service.emailVerification.Update(ctx, tx, emailVerification)
		if err != nil {
			return model.UserRegisterResponse{}, err
		}
	}

	return response, nil
}
