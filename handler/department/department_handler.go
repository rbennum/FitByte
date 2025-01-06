package department_handler

import (
	"net/http"

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
// @Summary This API will create a new department
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
		h.logger.Warn(err.Error(), helper.FunctionCaller("departmentHandler.Create"), input)
		ctx.JSON(http.StatusUnprocessableEntity, helper.NewResponse(nil, err))
		return
	}

	// response, err := h.service.CreateDepartment(*input)
	// TODO: check whether the err above is nil or not
	// ctx.JSON(http.StatusCreated, handler.NewResponse(response, nil))
}

// List all available departments
// @Tags department
// @Summary This API will fetch a list of all departments
// @Description List all available departments
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response{data=helper.Response} "Created"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/department [GET]
func (h *handler) GetAll(ctx *gin.Context) {
	// TODO: get available params
	// TODO: validate whether the params are valid; if not, refer to default values
	// TODO: call service
	// TODO: return JSON
}

// List all available departments
// @Tags department
// @Summary This API will fetch a list of all departments
// @Description List all available departments
// @Accept json
// @Produce json
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

// List all available departments
// @Tags department
// @Summary This API will fetch a list of all departments
// @Description List all available departments
// @Accept json
// @Produce json
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
