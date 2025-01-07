package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/di"
	authHandler "github.com/levensspel/go-gin-template/handler/auth"
	fileHandler "github.com/levensspel/go-gin-template/handler/file"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	fileRepository "github.com/levensspel/go-gin-template/repository/file"
	fileService "github.com/levensspel/go-gin-template/service/file"
	userService "github.com/levensspel/go-gin-template/service/user"
	"github.com/samber/do/v2"

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

	fileRepo := fileRepository.NewFileRepository(db)

	userService := do.MustInvoke[userService.UserService](di.Injector)
	fileService := fileService.NewFileService(fileRepo, logger)

	userHdlr := userHandler.NewUserHandler(userService, logger)
	authHandler := authHandler.NewHandler(userService, logger)
	fileHandler := fileHandler.NewHandler(fileService, logger)

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

		file := controllers.Group("/file")
		{
			file.POST("", middleware.Authorization, fileHandler.Upload)
		}

		user := controllers.Group("/user")
		{
			user.GET("", middleware.Authorization, userHdlr.GetProfile)
			user.PUT("", middleware.Authorization, userHdlr.Update)
			user.DELETE("", middleware.Authorization, userHdlr.Delete)
		}
		// tambah route lainnya disini
	}

}
