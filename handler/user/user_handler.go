package userHandler

import (
	"net/http"

	"github.com/TimDebug/FitByte/helper"
	"github.com/TimDebug/FitByte/logger"
	"github.com/TimDebug/FitByte/middleware"
	service "github.com/TimDebug/FitByte/service/user"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type UserHandler interface {
	Get(ctx *gin.Context)
}

type handler struct {
	service service.UserService
	logger  logger.Logger
}

func NewUserHandler(service service.UserService, logger logger.Logger) UserHandler {
	return &handler{service: service, logger: logger}
}

func NewUserHandlerInject(i do.Injector) (UserHandler, error) {
	_service := do.MustInvoke[service.UserService](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return NewUserHandler(_service, &_logger), nil
}

// Get Profile user
// @Tags users
// @Summary Get Profile User
// @Description Get Profile User
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer + user token"
// @Success 200 {object} helper.Response{data=helper.Response} "OK"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorization"
// @Router /v1/user [GET]
func (h handler) Get(ctx *gin.Context) {
	id, err := middleware.GetUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}

	response, err := h.service.GetProfile(ctx, id)
	if err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("UserHandler.GetProfile"), id)
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusOK, response)
}
