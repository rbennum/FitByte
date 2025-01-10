package userService

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/levensspel/go-gin-template/cache"

	"github.com/levensspel/go-gin-template/auth"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/entity"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/user"
	"github.com/levensspel/go-gin-template/validation"
	"github.com/samber/do/v2"
	"golang.org/x/crypto/bcrypt"
)

const IsUserReadHeavy = true // Caching is suitable for read heavy operations

type IUserService interface {
	RegisterUser(input dto.UserRequestPayload) (dto.ResponseRegister, error)
	Login(input dto.UserRequestPayload) (dto.ResponseLogin, error)
	Update(input dto.RequestRegister) (dto.Response, error)
	DeleteByID(id string) error
	GetProfile(managerid string) (*dto.ResposneGetProfile, error)
	UpdateProfile(managerid string, input dto.RequestUpdateProfile) (*dto.RequestUpdateProfile, error)
}

type UserService struct {
	userRepo repositories.UserRepository
	logger   logger.LogHandler
}

func NewUserService(
	userRepo repositories.UserRepository,
	logger logger.LogHandler,
) UserService {
	return UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func NewUserServiceInject(i do.Injector) (UserService, error) {
	_userRepo := do.MustInvoke[repositories.UserRepository](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return NewUserService(_userRepo, _logger), nil
}

func (s *UserService) RegisterUser(input dto.UserRequestPayload) (dto.ResponseRegister, error) {
	err := validation.ValidateUserCreate(input, s.userRepo)
	if err != nil {
		return dto.ResponseRegister{}, helper.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	_, found := cache.Get(fmt.Sprintf(cache.CacheAuthEmailToToken, input.Email))
	if found {
		return dto.ResponseRegister{}, fmt.Errorf("email %s is already in use", input.Email)
	}

	user := entity.User{}

	user.Email.String = input.Email
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		s.logger.Error(err.Error(), helper.GenerateFromPassword, passwordHash)
		return dto.ResponseRegister{}, err
	}
	user.Password = string(passwordHash)

	user.Id, err = s.userRepo.Create(context.Background(), user)

	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return dto.ResponseRegister{}, helper.ErrConflict
		} else {
			s.logger.Error(err.Error(), helper.UserServiceRegister, user)
			return dto.ResponseRegister{}, err
		}
	}

	jwtService := auth.NewJWTService()
	token, err := jwtService.GenerateToken(user.Id)

	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceRegister, err)
		return dto.ResponseRegister{}, err
	}

	response := dto.ResponseRegister{
		Email: user.Email.String,
		Token: token,
	}

	return response, nil
}

func (s *UserService) Login(input dto.UserRequestPayload) (dto.ResponseLogin, error) {
	err := validation.ValidateUserLogin(input)
	if err != nil {
		return dto.ResponseLogin{}, helper.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	// Check cache first
	cachedToken, found := cache.Get(fmt.Sprintf(cache.CacheAuthEmailToToken, input.Email))
	if found {
		return dto.ResponseLogin{
			Email: input.Email,
			Token: cachedToken,
		}, nil
	}

	//get user
	fmt.Printf("email %s", input.Email)
	user, err := s.userRepo.GetUserbyEmail(context.Background(), input.Email)
	if err != nil {
		s.logger.Error(err.Error(), helper.FunctionCaller("UserService.Login.GetUserbyEmail"), input)
		return dto.ResponseLogin{}, err
	}
	if len(user) == 0 {
		return dto.ResponseLogin{}, helper.NewErrorResponse(http.StatusNotFound, helper.ErrNotFound.Error())
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

	// Put to cache
	cache.Set(fmt.Sprintf(cache.CacheAuthEmailToToken, input.Email), token)

	response := dto.ResponseLogin{}
	response.Email = user[0].Email.String
	response.Token = token
	return response, nil
}

func (s *UserService) Update(input dto.RequestRegister) (dto.Response, error) {
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

	if IsUserReadHeavy {
		// Put to cache
		profileToCache := map[string]string{
			"email":           user.Email.String,
			"name":            user.Name.String,
			"userImageUri":    "", // TODO: Handle once retrieved in request
			"companyName":     "", // TODO: Handle once retrieved in request
			"companyImageUri": "", // TODO: Handle once retrieved in request
		}
		cache.SetAsMap(fmt.Sprintf(cache.CacheUserIdToProfile, user.Id), profileToCache)
	}

	response := dto.Response{}
	response.Id = input.Id
	response.Email = input.Email
	response.Username = input.Username

	return response, nil
}

func (s *UserService) DeleteByID(id string) error {
	cache.Delete(fmt.Sprintf(cache.CacheUserIdToProfile, id))

	err := s.userRepo.Delete(context.Background(), id)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return err
	}
	return err
}

// Get manager profile by their id
func (s *UserService) GetProfile(id string) (*dto.ResposneGetProfile, error) {

	// Get from cache
	cachedProfile, found := cache.GetAsMap(fmt.Sprintf(cache.CacheUserIdToProfile, id))
	if found {
		return &dto.ResposneGetProfile{
			Email:           cachedProfile["email"],
			Name:            cachedProfile["name"],
			UserImageUri:    cachedProfile["userImageUri"],
			CompanyName:     cachedProfile["companyName"],
			CompanyImageUri: cachedProfile["companyImageUri"],
		}, nil
	}

	profile, err := s.userRepo.GetProfile(context.Background(), id)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceGetProfile, err)
		return nil, err
	}

	if IsUserReadHeavy {
		// Put to cache
		profileToCache := map[string]string{
			"email":           profile.Email,
			"name":            profile.Name.String,
			"userImageUri":    profile.UserImageUri.String,
			"companyName":     profile.CompanyName.String,
			"companyImageUri": profile.CompanyImageUri.String,
		}
		cache.SetAsMap(fmt.Sprintf(cache.CacheUserIdToProfile, id), profileToCache)
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

// Update manager profile by their id
func (s *UserService) UpdateProfile(id string, req dto.RequestUpdateProfile) (*dto.RequestUpdateProfile, error) {
	profile, err := s.userRepo.GetProfile(context.Background(), id)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceGetProfile, err)
		return nil, err
	}

	if req.Email != nil && *req.Email != profile.Email {
		user, err := s.userRepo.GetUserbyEmail(context.Background(), *req.Email)
		if err != nil || len(user) != 0 {
			return nil, helper.ErrConflict
		}
	}

	if req.Email != nil {
		profile.Email = *req.Email
	}
	if req.Name != nil {
		profile.Name = ToNullString(req.Name)
	}
	if req.UserImageUri != nil {
		profile.UserImageUri = ToNullString(req.UserImageUri)
	}
	if req.CompanyName != nil {
		profile.CompanyName = ToNullString(req.CompanyName)
	}
	if req.CompanyImageUri != nil {
		profile.CompanyImageUri = ToNullString(req.CompanyImageUri)
	}

	s.userRepo.UpdateProfile(context.Background(), id, profile)

	result := dto.RequestUpdateProfile{
		Email:           &profile.Email,
		Name:            &profile.Name.String,
		UserImageUri:    &profile.UserImageUri.String,
		CompanyName:     &profile.CompanyName.String,
		CompanyImageUri: &profile.CompanyImageUri.String,
	}

	return &result, nil
}

func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}
