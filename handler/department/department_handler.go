package departmentHandler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	service "github.com/levensspel/go-gin-template/service/department"
	"github.com/samber/do/v2"
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

func NewInject(i do.Injector) (DepartmentHandler, error) {
	_service := do.MustInvoke[service.DepartmentService](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return New(_service, &_logger), nil
}

// Create a new department
// @Tags department
// @Summary Create a new department
// @Description Create a new department
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param data body dto.RequestDepartment true "data"
// @Success 201 {object} helper.Response{data=helper.Response} "Created"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/department [POST]
func (h *handler) Create(ctx *gin.Context) {
	managerID, err := middleware.GetIdUserFromContext(ctx)
	if err != nil {
		h.logger.Warn(err.Error(), helper.DepartmentHandlerCreate)
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	input := new(dto.RequestDepartment)
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn(err.Error(), helper.DepartmentHandlerCreate, input)
		ctx.JSON(http.StatusBadRequest, helper.NewResponse(nil, err)) // throw bad request immediately
		return
	}
	response, err := h.service.Create(managerID, *input)
	if err != nil {
		h.logger.Error(err.Error(), helper.DepartmentHandlerCreate)
		if errors.Is(err, helper.ErrBadRequest) {
			ctx.JSON(http.StatusBadRequest, helper.NewErrorResponse(http.StatusBadRequest, err.Error()))
			return
		} else if errors.Is(err, helper.ErrConflict) {
			ctx.JSON(http.StatusConflict, helper.NewErrorResponse(http.StatusConflict, err.Error()))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, helper.NewErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
	}
	ctx.JSON(http.StatusCreated, response)
}

// List all available departments
// @Tags department
// @Summary Fetch a list of all departments
// @Description List all available departments
// @Accept json
// @Produce json
// @Param limit query int false "limit query param"
// @Param offset query int false "offset query param"
// @Param name query string false "department name"
// @Param Authorization header string true "Bearer JWT token"
// @Success 200 {object} helper.Response{data=helper.Response} "Created"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/department [GET]
func (h *handler) GetAll(ctx *gin.Context) {
	managerID, err := middleware.GetIdUserFromContext(ctx)
	if err != nil {
		h.logger.Warn(err.Error(), helper.DepartmentHandlerCreate)
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	limit := h.getQueryInt(ctx, "limit", 5)
	offset := h.getQueryInt(ctx, "offset", 0)
	name := ctx.DefaultQuery("name", "")
	if name == "" {
		name = "%"
	} else {
		name = fmt.Sprintf("%%%s%%", name)
	}
	h.logger.Info(fmt.Sprintf("%d, %d, %s, %s", limit, offset, name, managerID), helper.DepartmentHandlerGetAll)
	input := dto.RequestDepartment{}
	input.DepartmentName = name
	input.Limit = limit
	input.Offset = offset
	response, err := h.service.GetAll(managerID, input)
	if err != nil {
		h.logger.Error(err.Error(), helper.FunctionCaller("handler.GetAll"))
		ctx.JSON(http.StatusBadGateway, helper.NewResponse(nil, err))
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// Update a single record of department
// @Tags department
// @Summary Update a single record of department
// @Description Update a single record of department
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer + user token"
// @Param data body dto.RequestDepartment true "data"
// @Param id path string true "department ID"
// @Success 200 {object} helper.Response{data=helper.Response} "Created"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/department/{id} [PATCH]
func (h *handler) Update(ctx *gin.Context) {
	managerID, err := middleware.GetIdUserFromContext(ctx)
	if err != nil {
		h.logger.Warn(err.Error(), helper.DepartmentHandlerPatch)
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	deptID := ctx.Param("id")
	input := new(dto.RequestDepartment)
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn(err.Error(), helper.DepartmentHandlerPatch, input)
		ctx.JSON(http.StatusBadRequest, helper.NewResponse(nil, err))
		return
	}
	response, err := h.service.Update(input.DepartmentName, deptID, managerID)
	if err != nil {
		h.logger.Error(err.Error(), helper.DepartmentHandlerPatch)
		ctx.JSON(helper.GetErrorStatusCode(err), err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// Delete a department
// @Tags department
// @Summary Delete a department
// @Description Delete a department
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer + user token"
// @Param id path string true "department ID"
// @Success 200 {object} helper.Response{data=helper.Response} "Created"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized"
// @Failure 500 {object} helper.Response{errors=helper.ErrorResponse} "Server Error"
// @Router /v1/department/{id} [DELETE]
func (h *handler) Delete(ctx *gin.Context) {
	managerID, err := middleware.GetIdUserFromContext(ctx)
	if err != nil {
		h.logger.Error(err.Error(), helper.DepartmentHandlerDelete)
		ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		return
	}
	deptID := ctx.Param("id")
	if deptID == "" {
		h.logger.Error(err.Error(), helper.DepartmentHandlerDelete)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	err = h.service.Delete(deptID, managerID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%s is not found", deptID)})
		} else if errors.Is(err, helper.ErrConflict) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "still contain employee(s)"})
		} else {
			ctx.JSON(helper.GetErrorStatusCode(err), helper.NewResponse(nil, err))
		}
		h.logger.Error(err.Error(), helper.DepartmentHandlerDelete, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": "the department has been successfully deleted"})
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
