package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/di"
	authHandler "github.com/levensspel/go-gin-template/handler/auth"
	departmentHandler "github.com/levensspel/go-gin-template/handler/department"
	fileHandler "github.com/levensspel/go-gin-template/handler/file"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	departmentRepository "github.com/levensspel/go-gin-template/repository/department"
	fileRepository "github.com/levensspel/go-gin-template/repository/file"
	departmentService "github.com/levensspel/go-gin-template/service/department"
	fileService "github.com/levensspel/go-gin-template/service/file"
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
	departmentRepo := departmentRepository.New(db)

	fileService := fileService.NewFileService(fileRepo, logger)
	departmentService := departmentService.New(departmentRepo, logger)

	userHandler := do.MustInvoke[userHandler.UserHandler](di.Injector)
	authHandler := do.MustInvoke[authHandler.AuthorizationHandler](di.Injector)

	fileHandler := fileHandler.NewHandler(fileService, logger)
	deptHandler := departmentHandler.New(departmentService, logger)

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
			user.GET("", middleware.Authorization, userHandler.GetProfile)
			user.PUT("", middleware.Authorization, userHandler.Update)
			user.DELETE("", middleware.Authorization, userHandler.Delete)
		}
		department := controllers.Group("/department")
		{
			department.POST("", middleware.Authorization, deptHandler.Create)
			department.GET("", middleware.Authorization, deptHandler.GetAll)
			department.PATCH("/:id", middleware.Authorization, deptHandler.Update)
			department.DELETE("/:id", middleware.Authorization, deptHandler.Delete)
		}
		// tambah route lainnya disini
	}

}
