package server

import (
	"github.com/TimDebug/FitByte/di"
	authHandler "github.com/TimDebug/FitByte/handler/auth"
	departmentHandler "github.com/TimDebug/FitByte/handler/department"
	employeeHandler "github.com/TimDebug/FitByte/handler/employee"
	fileHandler "github.com/TimDebug/FitByte/handler/file"
	userHandler "github.com/TimDebug/FitByte/handler/user"
	"github.com/TimDebug/FitByte/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"

	_ "github.com/TimDebug/FitByte/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(r *gin.Engine, db *pgxpool.Pool) {
	userHandler := do.MustInvoke[userHandler.UserHandler](di.Injector)
	authHandler := do.MustInvoke[authHandler.AuthorizationHandler](di.Injector)
	fileHandler := do.MustInvoke[fileHandler.FileHandler](di.Injector)
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
			user.PATCH("", middleware.Authorization, middleware.ContentType, userHandler.UpdateProfile)
			user.DELETE("", middleware.Authorization, userHandler.Delete)
		}
		department := controllers.Group("/department")
		{
			department.POST("", middleware.ContentType, middleware.Authorization, deptHandler.Create)
			department.GET("", middleware.Authorization, deptHandler.GetAll)
			department.PATCH("/:id", middleware.ContentType, middleware.Authorization, deptHandler.Update)
			department.DELETE("/:id", middleware.Authorization, deptHandler.Delete)
		}

		employee := controllers.Group("/employee")
		{
			employee.POST("", middleware.Authorization, middleware.ContentType, employeeHdlr.Create)
			employee.GET("", middleware.Authorization, employeeHdlr.GetAll)
			employee.PATCH(":identityNumber", middleware.Authorization, middleware.ContentType, employeeHdlr.Update)
			employee.DELETE(":identityNumber", middleware.Authorization, employeeHdlr.Delete)
		}
		// tambah route lainnya disini
	}

}
