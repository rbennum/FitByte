package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/auth"
	"github.com/levensspel/go-gin-template/helper"
)

func Authorization(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	fmt.Printf("Authorization Header Val: %s\b", authorizationHeader)
	if !strings.Contains(authorizationHeader, "Bearer") && !strings.Contains(authorizationHeader, "bearer") {
		c.JSON(http.StatusUnauthorized, helper.NewResponse(nil, errors.New("the request is allowed for logged in")))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	bearerToken := ""
	if strings.Contains(authorizationHeader, "Bearer") {
		bearerToken = strings.Replace(authorizationHeader, "Bearer ", "", -1)
	}
	if strings.Contains(authorizationHeader, "bearer") {
		bearerToken = strings.Replace(authorizationHeader, "bearer ", "", -1)
	}

	fmt.Printf("Token Val: %s\b", bearerToken)

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

func GetIdUserFromContext(ctx *gin.Context) (string, error) {
	id, ok := ctx.Value("user_id").(string)
	if !ok {
		log.Printf(`Failed get data user context %v`, ctx.Value("user_id"))
		return ``, fmt.Errorf("invalid user context")
	}
	return id, nil
}
