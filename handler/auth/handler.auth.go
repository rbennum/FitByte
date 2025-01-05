package authHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	service "github.com/levensspel/go-gin-template/service/user"
)

type AuthorizationHandler interface {
	Post(ctx *gin.Context)
}

type handler struct {
	service service.UserService
	logger  logger.Logger
}

func NewHandler(service service.UserService, logger logger.Logger) AuthorizationHandler {
	return &handler{service: service, logger: logger}
}

// Entry for authentication or create new user
// @Tags auth
// @Summary Entry for authentication or create new user
// @Description either create or login
// @Accept json
// @Produce json
// @Param data body dto.UserRequestPayload true "data"
// @Success 200 {object} helper.Response{data=helper.Response} "EXISTING"
// @Success 201 {object} helper.Response{data=helper.Response} "CREATED"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Router /v1/auth [POST]
func (h handler) Post(ctx *gin.Context) {
	input := new(dto.UserRequestPayload)

	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("UserHandler"), &input)
		ctx.JSON(http.StatusUnprocessableEntity, helper.NewResponse(nil, err))
		return
	}

	switch input.Action {
	case dto.Create:
		// do register/create new user
		modelState := make(map[string]string)
		if input.Email == "" {
			modelState["Email"] = "do not left Email empty"
		}
		if input.Password == "" {
			modelState["Password"] = "do not left Password emtpy"
		}
		if len(modelState) == 0 {
			response, err := h.service.RegisterUser(*input)
			if err != nil {
				ctx.JSON(
					helper.GetErrorStatusCode(err),
					helper.NewResponse(
						helper.ErrorResponse{
							Code:    helper.GetErrorStatusCode(err),
							Message: "Either username, email, or choosen password has been selected",
						},
						err,
					),
				)
				return
			}

			ctx.JSON(http.StatusCreated, helper.NewResponse(response, err))
		} else {
			ctx.JSON(http.StatusBadRequest, helper.NewResponse(modelState, nil))
		}

	case dto.Login:
		h.logger.Warn("Login method not implemented.", helper.FunctionCaller("Auth.Post"), input)
		// do login
	}
}

// Login user
// @Tags users
// @Summary Login user
// @Description Login user
// @Accept  json
// @Produce  json
// @Param data body dto.RequestLogin true "data"
// @Success 200 {object} helper.Response{data=dto.ResponseLogin} "OK"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 404 {object} helper.Response{errors=helper.ErrorResponse} "Record not found"
// @Router /v1/user/login [POST]
func (h handler) Login(ctx *gin.Context) {
	input := new(dto.RequestLogin)
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("Register"), input)
		ctx.JSON(http.StatusUnprocessableEntity, helper.NewResponse(helper.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "Please verify your input",
		}, err))
		return
	}
	response, err := h.service.Login(*input)

	if err != nil {
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusOK, helper.Response{Data: response, Error: err})
}
