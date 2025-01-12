package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/helper"
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
