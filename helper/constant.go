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

	ActivityHandlerGetAll FunctionCaller = "ActivityHandler.GetAll"
	ActivityServiceGetAll FunctionCaller = "ActivityService.GetAll"

	GenerateFromPassword FunctionCaller = "GenerateFromPassword"

	UserHandler FunctionCaller = "UserHandler"
)

var ErrorBadRequest = errors.New("invalid request format")
var ErrorNotFound = errors.New("data not found")
var ErrorInternalServerError = errors.New("internal server error")

var ErrorEmailRegistered = errors.New("email is already registered")
var ErrorUsernameRegistered = errors.New("username is already registered")
var ErrorInvalidLogin = NewErrorResponse(http.StatusBadRequest, "invalid email or password")

var WORK_DIR string
