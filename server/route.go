package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	authHandler "github.com/levensspel/go-gin-template/handler/auth"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	repositories "github.com/levensspel/go-gin-template/repository/user"
	userService "github.com/levensspel/go-gin-template/service/user"

	_ "github.com/levensspel/go-gin-template/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(r *gin.Engine, db *pgxpool.Pool) {
	logger := logger.NewlogHandler()

	// api := r.Group("/v1")
	// {
	// 	// untuk memanfaatkan api versioning, uncomment dan pakai ini
	// }

	userRepo := repositories.NewUserRepository(db)
	userService := userService.NewUserService(userRepo, logger)
	userHdlr := userHandler.NewUserHandler(userService, logger)
	authHandler := authHandler.NewHandler(userService, logger)

	swaggerRoute := r.Group("/")
	{
		//Route untuk Swagger
		swaggerRoute.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	controllers := r.Group("/v1")
	{
		auth := controllers.Group("/auth")
		{
			auth.POST("", authHandler.Post)
		}

		user := controllers.Group("/user")
		{
			user.POST("/login", userHdlr.Login)
			user.PUT("", middleware.Authorization, userHdlr.Update)
			user.DELETE("", middleware.Authorization, userHdlr.Delete)
		}
		// tambah route lainnya disini
	}

}
