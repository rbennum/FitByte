package server

import (
	"github.com/TimDebug/FitByte/di"
	userHandler "github.com/TimDebug/FitByte/handler/user"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"

	_ "github.com/TimDebug/FitByte/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(r *gin.Engine, db *pgxpool.Pool) {
	userHandler := do.MustInvoke[userHandler.UserHandler](di.Injector)

	swaggerRoute := r.Group("/")
	{
		//Route untuk Swagger
		swaggerRoute.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	controllers := r.Group("/v1")
	{
		user := controllers.Group("/user")
		{
			// user.GET("", middleware.Authorization, userHandler.Get)
			user.GET("", userHandler.Get)
		}
	}
}
