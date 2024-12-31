package server

import (
	"github.com/gin-gonic/gin"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	dbTrxRepository "github.com/levensspel/go-gin-template/repository/db_trx"
	userRepository "github.com/levensspel/go-gin-template/repository/user"
	userService "github.com/levensspel/go-gin-template/service/user"
	"gorm.io/gorm"
)

func NewRouter(r *gin.Engine, db *gorm.DB) {
	logger := logger.NewlogHandler()

	api := r.Group("/v1")
	dbTrxRepo := dbTrxRepository.NewDBTrxRepository(db)

	userRepo := userRepository.NewUserRepository(db)
	userSrv := userService.NewUserService(userRepo, dbTrxRepo, logger)
	userHdlr := userHandler.NewUserHandler(userSrv)

	userRoute := api.Group("/users")
	userRoute.POST("/register", userHdlr.Register)
	userRoute.POST("/login", userHdlr.Login)
	userRoute.PUT("", middleware.Authorization, userHdlr.Update)
	userRoute.DELETE("", middleware.Authorization, userHdlr.Delete)
}
