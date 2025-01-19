package server

import (
	"net/http"

	"github.com/TimDebug/FitByte/di"
	"github.com/TimDebug/FitByte/handler"
	"github.com/TimDebug/FitByte/helper"
	"github.com/TimDebug/FitByte/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"

	_ "github.com/TimDebug/FitByte/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(r *gin.Engine, db *pgxpool.Pool) {
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, helper.ErrNotFound)
	})

	swaggerRoute := r.Group("/")
	{
		//Route untuk Swagger
		swaggerRoute.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	userHandler := do.MustInvoke[handler.UserHandler](di.Injector)
	authHandler := do.MustInvoke[handler.AuthorizationHandler](di.Injector)
	activityHandler := do.MustInvoke[handler.ActivityHandler](di.Injector)

	controllers := r.Group("/v1")
	{
		controllers.POST("/login", authHandler.Login)
		controllers.POST("/register", authHandler.Register)
		user := controllers.Group("/user")
		{
			user.GET("", middleware.Authorization, userHandler.Get)
		}
		activity := controllers.Group("/activity")
		{
			activity.GET("", middleware.Authorization, activityHandler.GetAll)
		}
	}
}
