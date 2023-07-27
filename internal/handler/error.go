package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
