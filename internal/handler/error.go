package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrorInvalidUserId      = errors.New("invalid id")
	ErrorInvalidIdParameter = errors.New("invalid id parameter")
	ErrorInvalidAuthHeader  = errors.New("invalid auth header")
	ErrorEmptyAuthHeader    = errors.New("empty auth header")
	ErrorForbidden          = errors.New("forbidden")
)

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, err error) {
	logrus.Error(err.Error())
	c.AbortWithStatusJSON(statusCode, errorResponse{
		Error: err.Error(),
	})
}
