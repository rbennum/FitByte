package handler

import (
	"net/http"

	"github.com/TimDebug/FitByte/dto"
	"github.com/TimDebug/FitByte/helper"
	"github.com/TimDebug/FitByte/logger"
	"github.com/TimDebug/FitByte/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AuthorizationHandler interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authHandler struct {
	service service.UserService
	logger  logger.Logger
}

func NewHandler(service service.UserService, logger logger.Logger) AuthorizationHandler {
	return &authHandler{service: service, logger: logger}
}

func NewHandlerInject(i do.Injector) (AuthorizationHandler, error) {
	_service := do.MustInvoke[service.UserService](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return NewHandler(_service, &_logger), nil
}

// Login
// @Tags auth
// @Summary User Login
// @Description User Login
// @Param data body dto.UserRequestPayload true "data"
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response{data=helper.Response} "OK"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 404 {object} helper.Response{errors=helper.ErrorResponse} "Not Found"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/login [POST]
func (h authHandler) Login(ctx *gin.Context) {
	requestBody := new(dto.UserRequestPayload)

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("AuthHandler.Login"), &requestBody)
		ctx.JSON(http.StatusBadRequest, helper.NewResponse(nil, err))
		return
	}

	response, err := h.service.Login(ctx, requestBody)

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// Register
// @Tags auth
// @Summary User Register
// @Description User Register
// @Param data body dto.UserRequestPayload true "data"
// @Accept json
// @Produce json
// @Success 201 {object} helper.Response{data=helper.Response} "OK"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 409 {object} helper.Response{errors=helper.ErrorResponse} "Conflict"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/register [POST]
func (h authHandler) Register(ctx *gin.Context) {
	requestBody := new(dto.UserRequestPayload)

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("AuthHandler.Register"), &requestBody)
		ctx.JSON(http.StatusBadRequest, helper.NewResponse(nil, err))
		return
	}

	response, err := h.service.Register(ctx, requestBody)

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), err)
		return
	}
	ctx.JSON(http.StatusCreated, response)
}
