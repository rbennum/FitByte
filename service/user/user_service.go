package userService

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TimDebug/FitByte/cache"
	"github.com/gin-gonic/gin"

	"github.com/TimDebug/FitByte/dto"
	"github.com/TimDebug/FitByte/entity"
	"github.com/TimDebug/FitByte/helper"
	"github.com/TimDebug/FitByte/logger"
	repositories "github.com/TimDebug/FitByte/repository/user"
	"github.com/samber/do/v2"
)

const (
	isCachingBatchOfProfilesEnabled = true // Caching is suitable for read heavy operations
	cacheDefaultTtl                 = 5 * time.Minute
	maxMemoizedInvalidatedUserIds   = 100
)

type IUserService interface {
	GetProfile(ctx *gin.Context, managerid string) (*dto.ResponseGetProfile, error)
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

// Get manager profile by their id
func (s *UserService) GetProfile(ctx *gin.Context, id string) (*dto.ResponseGetProfile, error) {

	// Get from cache
	cachedProfile, found := cache.GetAsMap(fmt.Sprintf(cache.CacheUserIdToProfile, id))
	if found {
		height, err := strconv.Atoi(cachedProfile["height"])
		weight, err := strconv.Atoi(cachedProfile["weight"])
		if err != nil {
			height = 0
			weight = 0
		}
		return &dto.ResponseGetProfile{
			Preference: cachedProfile["preference"],
			WeightUnit: cachedProfile["weightUnit"],
			HeightUnit: cachedProfile["heightUnit"],
			Weight:     weight,
			Height:     height,
			Email:      cachedProfile["email"],
			Name:       cachedProfile["name"],
			ImageUri:   cachedProfile["imageUri"],
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
			if p.Id == &id {
				profile = &p
			}

			// Put to cache
			profileToCache := map[string]string{
				"email":    *profile.Email,
				"name":     *profile.Name,
				"imageUri": *profile.ImageUri,
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
			"email":    *profile.Email,
			"name":     *profile.Name,
			"imageUri": *profile.ImageUri,
		}
		cache.SetAsMap(fmt.Sprintf(cache.CacheUserIdToProfile, id), profileToCache)
	}

	result := dto.ResponseGetProfile{
		Email:    *profile.Email,
		Name:     *profile.Name,
		ImageUri: *profile.ImageUri,
	}
	return &result, nil
}
