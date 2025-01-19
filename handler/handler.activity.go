package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TimDebug/FitByte/helper"
	"github.com/TimDebug/FitByte/logger"
	"github.com/TimDebug/FitByte/middleware"
	"github.com/TimDebug/FitByte/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type ActivityType string

const (
	Walking    ActivityType = "Walking"
	Yoga       ActivityType = "Yoga"
	Stretching ActivityType = "Stretching"
	Cycling    ActivityType = "Cycling"
	Swimming   ActivityType = "Swimming"
	Dancing    ActivityType = "Dancing"
	Hiking     ActivityType = "Hiking"
	Running    ActivityType = "Running"
	HIIT       ActivityType = "HIIT"
	JumpRope   ActivityType = "JumpRope"
)

func GetCaloriesPerMinute(activityType ActivityType) (float64, error) {
	switch activityType {
	case Walking, Yoga, Stretching:
		return 4.0, nil
	case Cycling, Swimming, Dancing:
		return 8.0, nil
	case Hiking, Running, HIIT, JumpRope:
		return 10.0, nil
	default:
		return 0, helper.ErrBadRequest
	}
}

func IsValidActivityType(activityType ActivityType) bool {
	_, err := GetCaloriesPerMinute(activityType)
	return err == nil
}

type ActivityHandler struct {
	service service.ActivityService
	logger  logger.Logger
}

func NewActivityHandler(service service.ActivityService, logger logger.Logger) *ActivityHandler {
	return &ActivityHandler{service: service, logger: logger}
}

func NewActivityHandlerInject(i do.Injector) (ActivityHandler, error) {
	_service := do.MustInvoke[service.ActivityService](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return *NewActivityHandler(_service, &_logger), nil
}

// List all available activities
// @Tags activity
// @Summary Fetch a list of all activities
// @Description List all available activities
// @Accept json
// @Produce json
// @Param limit query int false "limit query param"
// @Param offset query int false "offset query param"
// @Param name query string false "activity name"
// @Param doneAtFrom query string false "done at from in ISO date"
// @Param doneAtTo query string false "done at from in ISO date"
// @Param caloriesBurnedMin query int false "calories burned minimum"
// @Param caloriesBurnedMax query int false "calories burned maximum"
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} helper.Response{data=helper.Response} "OK"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/activity [GET]
func (a *ActivityHandler) GetAll(ctx *gin.Context) {
	id, err := middleware.GetUserIdFromContext(ctx)
	if err != nil {
		a.logger.Warn(err.Error(), helper.ActivityHandlerGetAll)
		ctx.JSON(helper.GetErrorStatusCode(err), err)
		return
	}

	params := make(map[string]string)
	params["id"] = id
	params["activityType"] = ctx.DefaultQuery("activityType", "")
	params["doneAtFrom"] = ctx.DefaultQuery("doneAtFrom", "")
	params["doneAtTo"] = ctx.DefaultQuery("doneAtTo", "")
	params["caloriesBurnedMin"] = ctx.DefaultQuery("caloriesBurnedMin", "")
	params["caloriesBurnedMax"] = ctx.DefaultQuery("caloriesBurnedMax", "")
	response, err := a.service.GetAll(ctx, buildQueryParams(ctx, params))
	if err != nil {
		a.logger.Error(err.Error(), helper.ActivityHandlerGetAll)
		ctx.JSON(helper.GetErrorStatusCode(err), err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func getQueryInt(ctx *gin.Context, key string, defaultValue int) int {
	value, exists := ctx.GetQuery(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func buildQueryParams(ctx *gin.Context, params map[string]string) []interface{} {
	limit := getQueryInt(ctx, "limit", 5)
	offset := getQueryInt(ctx, "offset", 0)
	activityType := params["activityType"]
	doneAtFrom := params["doneAtFrom"]
	doneAtTo := params["doneAtTo"]
	caloriesBurnedMin := getQueryInt(ctx, "caloriesBurnedMin", 0)
	caloriesBurnedMax := getQueryInt(ctx, "caloriesBurnedMax", 0)
	id := params["id"]

	queryArgs := []interface{}{id, nil, nil, nil, nil, nil, limit, offset}

	// validate activityType
	if activityType != "" && IsValidActivityType(ActivityType(activityType)) {
		queryArgs[1] = activityType
	}

	// validate doneAtFrom
	if doneAtFrom != "" {
		if parsedDate, err := time.Parse(time.RFC3339, doneAtFrom); err == nil {
			queryArgs[2] = parsedDate
		}
	}

	// validate doneAtTo
	if doneAtTo != "" {
		if parsedDate, err := time.Parse(time.RFC3339, doneAtTo); err == nil {
			queryArgs[3] = parsedDate
		}
	}

	// validate caloriesBurnedMin
	if caloriesBurnedMin > 0 {
		queryArgs[4] = caloriesBurnedMin
	}

	// validate caloriesBurnedMax
	if caloriesBurnedMax > 0 {
		queryArgs[5] = caloriesBurnedMax
	}

	return queryArgs
}
