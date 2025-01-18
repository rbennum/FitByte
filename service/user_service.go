package service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TimDebug/FitByte/auth"
	"github.com/TimDebug/FitByte/cache"
	"github.com/TimDebug/FitByte/repository"
	"github.com/TimDebug/FitByte/validation"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/TimDebug/FitByte/dto"
	"github.com/TimDebug/FitByte/entity"
	"github.com/TimDebug/FitByte/helper"
	"github.com/TimDebug/FitByte/logger"
	"github.com/samber/do/v2"
)

const (
	isCachingBatchOfProfilesEnabled = true // Caching is suitable for read heavy operations
	cacheDefaultTtl                 = 5 * time.Minute
	maxMemoizedInvalidatedUserIds   = 100
)

type UserService struct {
	userRepo repository.UserRepository
	logger   logger.LogHandler
}

func NewUserService(
	userRepo repository.UserRepository,
	logger logger.LogHandler,
) UserService {
	return UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func NewUserServiceInject(i do.Injector) (UserService, error) {
	_userRepo := do.MustInvoke[repository.UserRepository](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return NewUserService(_userRepo, _logger), nil
}

func (s *UserService) Login(ctx *gin.Context, body *dto.UserRequestPayload) (*dto.ResponseAuth, error) {
	err := validation.ValidateUserCreate(*body)
	if err != nil {
		return &dto.ResponseAuth{}, helper.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	users, err := s.userRepo.Login(ctx, body)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceLogin)
		return nil, err
	}
	if len(users) == 0 {
		return nil, helper.NewErrorResponse(http.StatusNotFound, "Not found")
	}

	passwordHash := users[0].PasswordHash
	err = bcrypt.CompareHashAndPassword([]byte(*passwordHash), []byte(body.Password))
	if err != nil {
		return nil, helper.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	cachedToken, found := cache.Get(fmt.Sprintf(cache.CacheAuthEmailToToken, body.Email))
	if found {
		return &dto.ResponseAuth{
			Email: body.Email,
			Token: cachedToken,
		}, nil
	}

	jwtService := auth.NewJWTService()
	token, err := jwtService.GenerateToken(*users[0].Id)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceRegister, err)
		return nil, helper.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}
	return &dto.ResponseAuth{Email: body.Email, Token: token}, nil
}

func (s *UserService) Register(ctx *gin.Context, body *dto.UserRequestPayload) (*dto.ResponseAuth, error) {
	err := validation.ValidateUserCreate(*body)
	if err != nil {
		return &dto.ResponseAuth{}, helper.NewErrorResponse(http.StatusBadRequest, err.Error())
	}

	_, found := cache.Get(fmt.Sprintf(cache.CacheAuthEmailToToken, body.Email))
	if found {
		return &dto.ResponseAuth{}, helper.ErrConflict
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)
	if err != nil {
		s.logger.Error(err.Error(), helper.GenerateFromPassword, passwordHash)
		return &dto.ResponseAuth{}, err
	}

	user := entity.User{}
	user.Email = &body.Email
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = user.CreatedAt
	password := string(passwordHash)
	user.PasswordHash = &password
	userId, err := s.userRepo.Register(ctx, &user)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceRegister)
		if strings.Contains(err.Error(), "23505") {
			return nil, helper.ErrConflict
		}
		return nil, helper.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}

	jwtService := auth.NewJWTService()
	token, err := jwtService.GenerateToken(userId)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceRegister, err)
		return nil, helper.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}

	cache.Set(fmt.Sprintf(cache.CacheAuthEmailToToken, body.Email), token)
	appendToInvalidatedUserIds(userId)
	return &dto.ResponseAuth{Email: body.Email, Token: token}, nil
}

// Get manager profile by their id
func (s *UserService) GetProfile(ctx *gin.Context, id string) (*dto.ResponseGetProfile, error) {
	profile, err := s.userRepo.GetProfile(ctx, id)
	if err != nil {
		s.logger.Error(err.Error(), helper.UserServiceGetProfile, err)
		return nil, err
	}

	result := dto.ResponseGetProfile{
		Email:    *profile.Email,
		Name:     profile.Name,
		ImageUri: profile.ImageUri,
	}
	return &result, nil
}

func getValue(cache map[string]string, key string, asInt bool) interface{} {
	val, exists := cache[key]
	if !exists {
		return nil
	}
	if asInt {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return nil
		}
		return &intVal
	}
	return &val
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
