package authHandler

import (
	"github.com/TimDebug/FitByte/logger"
	service "github.com/TimDebug/FitByte/service/user"
	"github.com/gin-gonic/gin"
)

type AuthorizationHandler interface {
	Post(ctx *gin.Context)
}

type handler struct {
	service service.UserService
	logger  logger.Logger
}
