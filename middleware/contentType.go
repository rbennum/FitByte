package middleware

import (
	"net/http"

	"github.com/TimDebug/FitByte/helper"
	"github.com/gin-gonic/gin"
)

func ContentType(ctx *gin.Context) {
	if ctx.ContentType() != "application/json" {
		ctx.JSON(
			http.StatusBadRequest,
			helper.NewErrorResponse(http.StatusBadRequest, "Content-Type must be application/json"),
		)
		ctx.Abort()
		return
	}
	ctx.Next()
}
