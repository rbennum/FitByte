package userService

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

const (
	isCachingBatchOfProfilesEnabled = true // Caching is suitable for read heavy operations
	cacheDefaultTtl                 = 5 * time.Minute
	maxMemoizedInvalidatedUserIds   = 100
)

type IUserService interface {
	RegisterUser(ctx *gin.Context, input dto.UserRequestPayload) (dto.ResponseRegister, error)
	Login(ctx *gin.Context, input dto.UserRequestPayload) (dto.ResponseLogin, error)
	Update(ctx *gin.Context, input dto.RequestRegister) (dto.Response, error)
	DeleteByID(ctx *gin.Context, id string) error
	GetProfile(ctx *gin.Context, managerid string) (*dto.ResponseGetProfile, error)
	UpdateProfile(ctx *gin.Context, managerid string, input dto.RequestUpdateProfile) (*dto.RequestUpdateProfile, error)
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

func (s *UserService) RegisterUser(ctx *gin.Context, input dto.UserRequestPayload) (dto.ResponseRegister, error) {
	err := validation.ValidateUserCreate(input, s.userRepo)
	if err != nil {
		return dto.ResponseRegister{}, helper.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	_, found := cache.Get(fmt.Sprintf(cache.CacheAuthEmailToToken, input.Email))
	if found {
		// return dto.ResponseRegister{}, fmt.Errorf("email %s is already in use", input.Email)
		return dto.ResponseRegister{}, helper.ErrConflict
	}

	user := entity.User{}

	user.Email.String = input.Email
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	prehashed := prehashPassword(input.Password)
	passwordHash, err := bcrypt.GenerateFromPassword(prehashed, bcrypt.MinCost)

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

	// Put to cache
	cache.Set(fmt.Sprintf(cache.CacheAuthEmailToToken, input.Email), token)
	appendToInvalidatedUserIds(user.Id)

	response := dto.ResponseRegister{
		Email: user.Email.String,
		Token: token,
	}

	return response, nil
}

func (s *UserService) Login(ctx *gin.Context, input dto.UserRequestPayload) (dto.ResponseLogin, error) {
	err := validation.ValidateUserLogin(input)
	if err != nil {
		return dto.ResponseLogin{}, helper.NewErrorResponse(http.StatusBadRequest, err.Error())
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
	prehashed := prehashPassword(input.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user[0].Password), prehashed)
	if err != nil {
		s.logger.Error(err.Error(), helper.FunctionCaller("UserService.Login.CompareHashAndPassword"), err)
		return dto.ResponseLogin{}, helper.ErrorInvalidLogin
	}

	// Get from cache
	cachedToken, found := cache.Get(fmt.Sprintf(cache.CacheAuthEmailToToken, input.Email))
	if found {
		return dto.ResponseLogin{
			Email: input.Email,
			Token: cachedToken,
		}, nil
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

func (s *UserService) Update(ctx *gin.Context, input dto.RequestRegister) (dto.Response, error) {
	prehashed := prehashPassword(input.Password)
	passwordHash, err := bcrypt.GenerateFromPassword(prehashed, bcrypt.MinCost)

	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return dto.Response{}, err
	}

	user := entity.User{}
	user.Id = input.Id
	user.Username.String = input.Username
	user.Email.String = input.Email
	user.Password = string(passwordHash)
	user.UpdatedAt = time.Now().Unix()
	err = s.userRepo.Update(context.Background(), user)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return dto.Response{}, err
	}

	// Invalidate cache
	cache.Delete(fmt.Sprintf(cache.CacheUserIdToProfile, input.Id))
	appendToInvalidatedUserIds(input.Id)

	response := dto.Response{}
	response.Id = input.Id
	response.Email = input.Email
	response.Username = input.Username

	return response, nil
}

func (s *UserService) DeleteByID(ctx *gin.Context, id string) error {
	// Invalidate cache
	cachedProfile, found := cache.GetAsMap(fmt.Sprintf(cache.CacheUserIdToProfile, id))
	if found {
		cache.Delete(fmt.Sprintf(cache.CacheAuthEmailToToken, cachedProfile["email"]))
		cache.Delete(fmt.Sprintf(cache.CacheUserIdToProfile, id))
	}

	err := s.userRepo.Delete(context.Background(), id)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceUpdate, err)
		return err
	}
	return err
}

// Get manager profile by their id
func (s *UserService) GetProfile(ctx *gin.Context, id string) (*dto.ResponseGetProfile, error) {

	// Get from cache
	cachedProfile, found := cache.GetAsMap(fmt.Sprintf(cache.CacheUserIdToProfile, id))
	if found {
		return &dto.ResponseGetProfile{
			Email:           cachedProfile["email"],
			Name:            cachedProfile["name"],
			UserImageUri:    cachedProfile["userImageUri"],
			CompanyName:     cachedProfile["companyName"],
			CompanyImageUri: cachedProfile["companyImageUri"],
		}, nil
	}

	var profile *entity.GetProfile
	var err error

	if isCachingBatchOfProfilesEnabled {
		invalidatedUserIds := make([]string, 0)
		v, found := cache.Get(cache.CacheInvalidatedUserIds)
		if found {
			invalidatedUserIds = strings.Split(v, ",")
		}
		invalidatedUserIds = append(invalidatedUserIds, id)

		batchOfProfiles, err := s.userRepo.GetBatchOfProfiles(context.Background(), invalidatedUserIds)
		if err != nil {
			s.logger.Error(err.Error(), helper.UserServiceGetProfile, err)
			return nil, err
		}

		for _, p := range batchOfProfiles {
			if p.ManagerId == id {
				profile = &p
			}

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

		cache.Delete(cache.CacheInvalidatedUserIds)

	} else {
		profile, err = s.userRepo.GetProfile(context.Background(), id)
		if err != nil {
			s.logger.Error(err.Error(), helper.UserServiceGetProfile, err)
			return nil, err
		}

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

	result := dto.ResponseGetProfile{
		Email:           profile.Email,
		Name:            profile.Name.String,
		UserImageUri:    profile.UserImageUri.String,
		CompanyName:     profile.CompanyName.String,
		CompanyImageUri: profile.CompanyImageUri.String,
	}
	return &result, nil
}

// Update manager profile by their id
func (s *UserService) UpdateProfile(ctx *gin.Context, id string, req dto.RequestUpdateProfile) (*dto.RequestUpdateProfile, error) {
	var profile *entity.GetProfile
	var err error

	// Get from cache
	cachedProfile, found := cache.GetAsMap(fmt.Sprintf(cache.CacheUserIdToProfile, id))
	if found {
		profile = &entity.GetProfile{
			Email:           cachedProfile["email"],
			Name:            sql.NullString{String: cachedProfile["name"], Valid: cachedProfile["name"] != ""},
			UserImageUri:    sql.NullString{String: cachedProfile["userImageUri"], Valid: cachedProfile["userImageUri"] != ""},
			CompanyName:     sql.NullString{String: cachedProfile["companyName"], Valid: cachedProfile["companyName"] != ""},
			CompanyImageUri: sql.NullString{String: cachedProfile["companyImageUri"], Valid: cachedProfile["companyImageUri"] != ""},
		}

	} else {
		profile, err = s.userRepo.GetProfile(context.Background(), id)
		if err != nil {
			s.logger.Error(err.Error(), helper.UserServiceGetProfile, err)
			return nil, err
		}
	}

	if req.Email != nil && *req.Email != profile.Email {
		// Check cache
		_, found := cache.Get(fmt.Sprintf(cache.CacheAuthEmailToToken, *req.Email))
		if found {
			return nil, helper.ErrConflict

		} else {
			user, err := s.userRepo.GetUserbyEmail(context.Background(), *req.Email)
			if err != nil || len(user) != 0 {
				return nil, helper.ErrConflict
			}
		}

		// Email is updated, then invalidate old auth email cache
		cache.Delete(fmt.Sprintf(cache.CacheAuthEmailToToken, profile.Email))
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

	// Invalidate profile cache
	cache.Delete(fmt.Sprintf(cache.CacheUserIdToProfile, id))
	appendToInvalidatedUserIds(id)

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

func appendToInvalidatedUserIds(id string) {
	invalidatedUserIds := make([]string, 0)
	v, found := cache.Get(cache.CacheInvalidatedUserIds)
	if found {
		invalidatedUserIds = strings.Split(v, ",")
	}

	if len(invalidatedUserIds) >= maxMemoizedInvalidatedUserIds {
		// Rooms for improvement: use circular buffer data structure
		// To replace oldest data with new data
		return
	}

	invalidatedUserIds = append(invalidatedUserIds, id)
	cache.Set(cache.CacheInvalidatedUserIds, strings.Join(invalidatedUserIds, ","))
}

func prehashPassword(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}
