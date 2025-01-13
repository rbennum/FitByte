package userHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	service "github.com/levensspel/go-gin-template/service/user"
	"github.com/levensspel/go-gin-template/validation"
	"github.com/samber/do/v2"
)

type UserHandler interface {
	Update(ctx *gin.Context) // Deprecated, replaced by UpdateProfile
	Delete(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
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

// Update user
// @Tags users
// @Summary Update user
// @Description Update user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer + user token"
// @Param data body dto.UserRequestUpdate true "data"
// @Success 200 {object} helper.Response{data=helper.Response} "OK"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorization"
// @Router /v1/user [PUT]
func (h handler) Update(ctx *gin.Context) {
	input := new(dto.RequestRegister)

	if err := ctx.ShouldBind(input); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, helper.NewResponse(nil, err))
		return
	}
	id := ctx.MustGet("user_id")
	input.Id = id.(string)
	response, err := h.service.Update(ctx, *input)

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusOK, helper.NewResponse(response, nil))
}

// DeleteByID user
// @Tags users
// @Summary Delete user
// @Description Delete user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer + user token"
// @Success 200 {object} helper.Response{data=helper.Response} "OK"
// @Failure 404 {object} helper.Response{errors=helper.ErrorResponse} "Not Found"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorization"
// @Router /v1/user [DELETE]
func (h handler) Delete(ctx *gin.Context) {
	id := ctx.MustGet("user_id")

	err := h.service.DeleteByID(ctx, id.(string))

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}

	message := map[string]interface{}{"message": "your account has been successfully deleted"}
	ctx.JSON(http.StatusOK, helper.NewResponse(message, nil))
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
func (h handler) GetProfile(ctx *gin.Context) {
	id, err := middleware.GetIdUserFromContext(ctx)
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

// Update profile
// @Tags users
// @Summary Update profile
// @Description Update profile
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer + user token"
// @Param data body dto.RequestUpdateProfile true "data"
// @Success 200 {object} helper.Response{data=helper.Response} "OK"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorization"
// @Router /v1/user [PATCH]
func (h handler) UpdateProfile(ctx *gin.Context) {
	id, err := middleware.GetIdUserFromContext(ctx)
	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}

	req := new(dto.RequestUpdateProfile)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(err.Error(), helper.UserHandler, &req)
		ctx.JSON(http.StatusBadRequest, helper.NewResponse(nil, err))
		return
	}

	err = validation.ValidateUpdateProfile(*req)
	if err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("UserHandler.UpdateProfile"), &req)
		ctx.JSON(http.StatusBadRequest, helper.NewResponse(nil, err))
		return
	}

	response, err := h.service.UpdateProfile(ctx, id, *req)
	if err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("UserHandler.UpdateProfile"), id, &req)
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusOK, response)
}
