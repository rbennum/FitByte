package helper

import (
	"errors"
	"net/http"
)

type FunctionCaller string

const (
	UserRepoCreate FunctionCaller = "userRepo.Create"
	DbTrxRepoBegin FunctionCaller = "dbTrxRepo.Begin"

	UserServiceRegister   FunctionCaller = "userService.RegisterUser"
	UserServiceLogin      FunctionCaller = "userService.Login"
	UserServiceUpdate     FunctionCaller = "userService.Update"
	UserServiceDeleteByID FunctionCaller = "userService.DeleteById"
	UserServiceGetProfile FunctionCaller = "userService.GetProfile"

	EmployeeHandlerCreate       FunctionCaller = "EmployeeHandler.Create"
	EmployeeHandlerGetEmployees FunctionCaller = "EmployeeHandler.GetEmployees"
	EmployeeHandlerUpdate       FunctionCaller = "EmployeeHandler.Update"
	EmployeeHandlerDelete       FunctionCaller = "EmployeeHandler.Delete"

	EmployeeServiceCreate FunctionCaller = "employeeService.Create"
	EmployeeServiceGet    FunctionCaller = "employeeService.Get"
	EmployeeServiceUpdate FunctionCaller = "employeeService.Update"
	EmployeeServiceDelete FunctionCaller = "employeeService.Delete"

	GenerateFromPassword FunctionCaller = "GenerateFromPassword"

	UserHandler FunctionCaller = "UserHandler"

	DepartmentHandlerCreate FunctionCaller = "DepartmentHandler.Create"
	DepartmentHandlerGetAll FunctionCaller = "DepartmentHandler.GetAll"
	DepartmentHandlerPatch  FunctionCaller = "DepartmentHandler.Patch"
	DepartmentHandlerDelete FunctionCaller = "DepartmentHandler.Delete"

	DepartmentServiceCreate FunctionCaller = "DepartmentService.Create"
	DepartmentServiceGetAll FunctionCaller = "DepartmentService.GetAll"
	DepartmentServicePatch  FunctionCaller = "DepartmentService.Patch"
	DepartmentServiceDelete FunctionCaller = "DepartmentService.Delete"
)

var ErrorBadRequest = errors.New("invalid request format")
var ErrorNotFound = errors.New("data not found")
var ErrorInternalServerError = errors.New("internal server error")

var ErrorEmailRegistered = errors.New("email is already registered")
var ErrorUsernameRegistered = errors.New("username is already registered")
var ErrorInvalidLogin = NewErrorResponse(http.StatusBadRequest, "invalid email or password")

var WORK_DIR string
