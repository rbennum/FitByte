package departmentHandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	service "github.com/levensspel/go-gin-template/service/department"
)

type DepartmentHandler interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type handler struct {
	service service.DepartmentService
	logger  logger.Logger
}

func New(
	service service.DepartmentService,
	logger logger.Logger,
) DepartmentHandler {
	return &handler{service: service, logger: logger}
}

// Create a new department
// @Tags department
// @Summary Create a new department
// @Description Create a new department
// @Accept json
// @Produce json
// @Param data body dto.RequestDepartment true "data"
// @Success 201 {object} helper.Response{data=helper.Response} "Created"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/department [POST]
func (h *handler) Create(ctx *gin.Context) {
	input := new(dto.RequestDepartment)
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn(err.Error(), helper.DepartmentHandlerCreate, input)
		ctx.JSON(http.StatusUnprocessableEntity, helper.NewResponse(nil, err))
		return
	}
	response, err := h.service.Create(*input)
	if errors.Is(err, helper.ErrBadRequest) {
		h.logger.Error(err.Error(), helper.DepartmentHandlerCreate)
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	} else if err != nil {
		h.logger.Error(err.Error(), helper.DepartmentHandlerCreate, err)
		ctx.JSON(http.StatusInternalServerError, helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

// List all available departments
// @Tags department
// @Summary Fetch a list of all departments
// @Description List all available departments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer + user token"
// @Success 200 {object} helper.Response{data=helper.Response} "Created"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/department [GET]
func (h *handler) GetAll(ctx *gin.Context) {
	limit := h.getQueryInt(ctx, "limit", 5)
	offset := h.getQueryInt(ctx, "offset", 5)
	name := ctx.DefaultQuery("name", "")
	response, err := h.service.GetAll(name, limit, offset)
	if err != nil {
		h.logger.Error(err.Error(), helper.FunctionCaller("handler.GetAll"))
		ctx.JSON(http.StatusBadGateway, helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusAccepted, helper.NewResponse(response, nil))
}

// Update a single record of department
// @Tags department
// @Summary Update a single record of department
// @Description List all available departments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer + user token"
// @Param data body dto.RequestDepartment true "data"
// @Success 200 {object} helper.Response{data=helper.Response} "Created"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/department [PATCH]
func (h *handler) Update(ctx *gin.Context) {
	// TODO: get request params
	// TODO: validate whether the params are valid; if not, refer to default values
	// TODO: call service
	// TODO: return JSON
}

// Delete a department
// @Tags department
// @Summary Delete a department
// @Description List all available departments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer + user token"
// @Param data body dto.RequestDepartment true "data"
// @Success 200 {object} helper.Response{data=helper.Response} "Created"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/department [DELETE]
func (h *handler) Delete(ctx *gin.Context) {
	// TODO: get request params
	// TODO: validate whether the params are valid; if not, refer to default values
	// TODO: call service
	// TODO: return JSON
}

func (h *handler) getQueryInt(ctx *gin.Context, key string, defaultValue int) int {
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
