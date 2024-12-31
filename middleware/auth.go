package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/auth"
	"github.com/levensspel/go-gin-template/helper"
)

func Authorization(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		c.JSON(http.StatusUnauthorized, helper.NewResponse(nil, errors.New("the request is allowed for logged in")))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	bearerToken := strings.Replace(authorizationHeader, "Bearer ", "", -1)

	id, err := auth.ParseToken(bearerToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, helper.NewResponse(nil, err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("user_id", id)
	c.Next()
}
