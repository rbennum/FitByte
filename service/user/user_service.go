package userService

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/levensspel/go-gin-template/auth"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/entity"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/user"
	"github.com/levensspel/go-gin-template/validation"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(input dto.UserRequestPayload) (dto.ResponseRegister, error)
	Login(input dto.UserRequestPayload) (dto.ResponseLogin, error)
	Update(input dto.RequestRegister) (dto.Response, error)
	DeleteByID(id string) error
	GetProfile(managerid string) (*dto.ResposneGetProfile, error)
}

type service struct {
	userRepo repositories.UserRepository
	logger   logger.Logger
}

func NewUserService(
	userRepo repositories.UserRepository,
	logger logger.Logger,
) UserService {
	return &service{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *service) RegisterUser(input dto.UserRequestPayload) (dto.ResponseRegister, error) {
	err := validation.ValidateUserCreate(input, s.userRepo)

	if err != nil {
		return dto.ResponseRegister{}, err
	}

	user := entity.User{}

	user.Id = uuid.New().String()
	user.Email.String = input.Email
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		s.logger.Error(err.Error(), helper.GenerateFromPassword, passwordHash)
		return dto.ResponseRegister{}, err
	}
	user.Password = string(passwordHash)

	err = s.userRepo.Create(context.Background(), user)

	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return dto.ResponseRegister{}, helper.ErrConflict
		} else {
			s.logger.Error(err.Error(), helper.UserServiceRegister, user)
			return dto.ResponseRegister{}, err
		}
	}

	response := dto.ResponseRegister{
		Email: user.Email.String,
		Token: user.Id,
	}

	return response, nil
}

func (s *service) Login(input dto.UserRequestPayload) (dto.ResponseLogin, error) {
	err := validation.ValidateUserLogin(input)
	if err != nil {
		return dto.ResponseLogin{}, err
	}

	//get user
	user, err := s.userRepo.GetUserbyEmail(context.Background(), input.Email)
	if err != nil {
		s.logger.Error(err.Error(), helper.FunctionCaller("UserService.Login.GetUserbyEmail"), input)
		return dto.ResponseLogin{}, err
	}
	if len(user) == 0 {
		return dto.ResponseLogin{}, helper.ErrNotFound
	}

	// password compared
	err = bcrypt.CompareHashAndPassword([]byte(user[0].Password), []byte(input.Password))
	if err != nil {
		s.logger.Error(err.Error(), helper.FunctionCaller("UserService.Login.CompareHashAndPassword"), err)
		return dto.ResponseLogin{}, helper.ErrorInvalidLogin
	}

	jwtService := auth.NewJWTService()
	token, err := jwtService.GenerateToken(user[0].Id)

	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceLogin, err)
		return dto.ResponseLogin{}, err
	}

	response := dto.ResponseLogin{}
	response.Email = user[0].Email.String
	response.Token = token
	return response, nil
}

func (s *service) Update(input dto.RequestRegister) (dto.Response, error) {
	user := entity.User{}
	user.Id = input.Id
	user.Username.String = input.Username
	user.Email.String = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return dto.Response{}, err
	}

	user.Password = string(passwordHash)
	user.UpdatedAt = time.Now().Unix()
	err = s.userRepo.Update(context.Background(), user)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return dto.Response{}, err
	}
	response := dto.Response{}
	response.Id = input.Id
	response.Email = input.Email
	response.Username = input.Username

	return response, nil
}

func (s *service) DeleteByID(id string) error {
	err := s.userRepo.Delete(context.Background(), id)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return err
	}
	return err
}

// Get manager profile by their id
func (s *service) GetProfile(id string) (*dto.ResposneGetProfile, error) {
	profile, err := s.userRepo.GetProfile(context.Background(), id)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceGetProfile, err)
		return nil, err
	}
	result := dto.ResposneGetProfile{
		Email:           profile.Email,
		Name:            profile.Name.String,
		UserImageUri:    profile.UserImageUri.String,
		CompanyName:     profile.CompanyName.String,
		CompanyImageUri: profile.CompanyImageUri.String,
	}
	return &result, nil
}
