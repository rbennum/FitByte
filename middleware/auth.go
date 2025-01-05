package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/auth"
	"github.com/levensspel/go-gin-template/helper"
)

func Authorization(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	fmt.Printf("Authorization Header: %s\n", authorizationHeader)
	if !strings.Contains(authorizationHeader, "Bearer") {
		c.JSON(http.StatusUnauthorized, helper.NewResponse(nil, errors.New("the request is allowed for logged in")))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	bearerToken := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	fmt.Printf("Auth Token: %s\n", bearerToken)
	id, err := auth.ParseToken(bearerToken)
	fmt.Printf("ID Parsed Token: %s\n", id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, helper.NewResponse(nil, err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("user_id", id)
	c.Next()
}
