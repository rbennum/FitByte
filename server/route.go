package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/di"
	authHandler "github.com/levensspel/go-gin-template/handler/auth"
	departmentHandler "github.com/levensspel/go-gin-template/handler/department"
	employeeHandler "github.com/levensspel/go-gin-template/handler/employee"
	fileHandler "github.com/levensspel/go-gin-template/handler/file"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/logger"
	"github.com/levensspel/go-gin-template/middleware"
	fileRepository "github.com/levensspel/go-gin-template/repository/file"
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
	fileService := fileService.NewFileService(fileRepo, logger)

	userHandler := do.MustInvoke[userHandler.UserHandler](di.Injector)
	authHandler := do.MustInvoke[authHandler.AuthorizationHandler](di.Injector)
	fileHandler := fileHandler.NewHandler(fileService, logger)
	deptHandler := do.MustInvoke[departmentHandler.DepartmentHandler](di.Injector)
	employeeHdlr := do.MustInvoke[employeeHandler.EmployeeHandler](di.Injector)

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

		employee := controllers.Group("/employee")
		{
			employee.POST("", middleware.Authorization, employeeHdlr.Create)
			employee.GET("", middleware.Authorization, employeeHdlr.GetAll)
			employee.DELETE(":identityNumber", middleware.Authorization, employeeHdlr.Delete)
		}
		// tambah route lainnya disini
	}

}
