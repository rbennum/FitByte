package user_service

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/levensspel/go-gin-template/auth"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/entity"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	dbTrxRepository "github.com/levensspel/go-gin-template/repository/db_trx"
	userRepository "github.com/levensspel/go-gin-template/repository/user"
	"github.com/levensspel/go-gin-template/validation"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(input dto.RequestRegister) (dto.ResponseRegister, error)
	Login(input dto.RequestLogin) (dto.ResponseLogin, error)
	Update(input dto.RequestRegister) (dto.Response, error)
	DeleteByID(input string) error
}

type service struct {
	userRepo  userRepository.UserRepository
	dbTrxRepo dbTrxRepository.DBTrxRepository
	logger    logger.Logger
}

func NewUserService(
	userRepo userRepository.UserRepository,
	dbTrxRepo dbTrxRepository.DBTrxRepository,
	logger logger.Logger,
) UserService {
	return &service{
		userRepo:  userRepo,
		dbTrxRepo: dbTrxRepo,
		logger:    logger,
	}
}

func (s *service) RegisterUser(input dto.RequestRegister) (dto.Response, error) {
	err := validation.ValidateUserCreate(input, s.userRepo)

	if err != nil {
		return dto.Response{}, err
	}

	user := entity.User{}

	user.Id = uuid.New().String()
	user.Username = input.Username
	user.Email = input.Email
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return dto.Response{}, err
	}

	user.Password = string(passwordHash)

	err = s.userRepo.Create(user, nil)
	if err != nil {
		return dto.Response{}, err
	}
	response := dto.Response{}
	copier.Copy(&response, &user)
	return response, nil
}

func (s *service) Login(input dto.RequestLogin) (dto.ResponseLogin, error) {
	err := validation.ValidateUserLogin(input)
	if err != nil {
		return dto.ResponseLogin{}, err
	}
	email := input.Email
	password := input.Password
	user, err := s.userRepo.FindByEmail(email, nil)

	if err != nil {
		return dto.ResponseLogin{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return dto.ResponseLogin{}, helper.ErrorInvalidLogin
	}

	jwtService := auth.NewJWTService()

	token, err := jwtService.GenerateToken(user.Id)

	if err != nil {
		return dto.ResponseLogin{}, err
	}

	response := dto.ResponseLogin{}
	response.Token = token
	return response, nil
}

func (s *service) Update(input dto.RequestRegister) (dto.Response, error) {
	user := entity.User{}
	user.Id = input.Id
	user.Username = input.Username
	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return dto.Response{}, err
	}

	user.Password = string(passwordHash)
	user.UpdatedAt = time.Now().Unix()
	updatedUser, err := s.userRepo.Update(user, nil)
	if err != nil {
		return dto.Response{}, err
	}
	response := dto.Response{}
	copier.Copy(&response, &updatedUser)

	return response, nil
}

func (s *service) DeleteByID(id string) error {
	return s.userRepo.DeleteByID(id, nil)
}
