package authHandler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	service "github.com/levensspel/go-gin-template/service/user"
	"github.com/samber/do/v2"
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

func NewHandlerInject(i do.Injector) (AuthorizationHandler, error) {
	_service := do.MustInvoke[service.UserService](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return NewHandler(_service, &_logger), nil
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
// @Failure 404 {object} helper.Response{errors=helper.ErrorResponse} "Not Found"
// @Failure 409 {object} helper.Response{errors=helper.ErrorResponse} "Conflict"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/auth [POST]
func (h handler) Post(ctx *gin.Context) {
	input := new(dto.UserRequestPayload)

	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("AuthHandler.Post"), &input)
		ctx.JSON(http.StatusBadRequest, helper.NewResponse(nil, err))
		return
	}

	switch strings.ToLower(input.Action) {
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
			h.logger.Info("Register", helper.FunctionCaller("AuthHander.Post"), input)
			response, err := h.service.RegisterUser(ctx, *input)
			h.logger.Info("After Register", helper.FunctionCaller("AuthHander.Post"), input)
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
			h.logger.Info("Created", helper.FunctionCaller("AuthHander.Post"), response)
			ctx.JSON(http.StatusCreated, response)
		} else {
			h.logger.Error("BadRequest", helper.FunctionCaller("AuthHander.Post"), modelState)
			ctx.JSON(http.StatusBadRequest, helper.NewResponse(modelState, nil))
		}
	case dto.Login:
		// do login
		fmt.Printf("input %s", *input)
		response, err := h.service.Login(ctx, *input)

		if err != nil {
			ctx.JSON(
				helper.GetErrorStatusCode(err),
				helper.NewResponse(
					helper.ErrorResponse{
						Code:    helper.GetErrorStatusCode(err),
						Message: helper.GetErrorMessage(err),
					},
					err,
				),
			)
			return
		}
		ctx.JSON(http.StatusOK, response)
	default:
		ctx.JSON(
			http.StatusBadRequest,
			helper.NewResponse(
				helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Action not found",
				},
				nil,
			),
		)
	}
}
