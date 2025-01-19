package service

import (
	"fmt"
	"net/http"

	"github.com/TimDebug/FitByte/dto"
	"github.com/TimDebug/FitByte/helper"
	"github.com/TimDebug/FitByte/logger"
	"github.com/TimDebug/FitByte/repository"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type ActivityService struct {
	repo   repository.ActivityRepository
	logger logger.LogHandler
}

func NewActivityService(
	repo repository.ActivityRepository,
	logger logger.LogHandler,
) ActivityService {
	return ActivityService{
		repo:   repo,
		logger: logger,
	}
}

func NewActivityServiceInject(i do.Injector) (ActivityService, error) {
	_repo := do.MustInvoke[repository.ActivityRepository](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return NewActivityService(_repo, _logger), nil
}

func (a *ActivityService) GetAll(ctx *gin.Context, queryArgs []interface{}) ([]dto.ResponseActivity, error) {
	a.logger.Info("param", helper.ActivityServiceGetAll, queryArgs...)
	rawActivities, err := a.repo.GetAll(ctx, queryArgs)
	if err != nil {
		a.logger.Error(err.Error(), helper.ActivityServiceGetAll, rawActivities)
		return nil, helper.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}

	returnedActivities := make([]dto.ResponseActivity, 0)
	for _, elem := range rawActivities {
		var activity dto.ResponseActivity
		activity.Id = *elem.ActivityId
		activity.ActivityType = *elem.ActivityType
		activity.DoneAt = fmt.Sprintf("%d", *elem.DoneAt)
		activity.DurationInMinutes = int(*elem.DurationInMinutes)
		activity.CaloriesBurned = int(*elem.CaloriesBurned)
		activity.CreatedAt = fmt.Sprintf("%d", *elem.CreatedAt)
		returnedActivities = append(returnedActivities, activity)
	}
	return returnedActivities, nil
}
