package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func validHeader(c *gin.Context) (string, bool) {
	header := c.Request.Header.Get("Authorization")

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, errors.New("empty auth header"))
		return "", false
	}

	parts := strings.Fields(header)

	if len(parts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, errors.New("invalid auth header"))
		return "", false
	}

	if parts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, errors.New("invalid auth header"))
		return "", false
	}

	if parts[1] == "" {
		newErrorResponse(c, http.StatusUnauthorized, errors.New("empty token"))
		return "", false
	}
	return parts[1], true
}

func (h *Handler) adminIdentity(c *gin.Context) {

	token, ok := validHeader(c)
	if !ok {
		return
	}

	if err := h.services.Admin.ParseToken(token); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
	}
}

func (h *Handler) trainerIdentity(c *gin.Context) {
	token, ok := validHeader(c)
	if !ok {
		return
	}

	if err := h.services.Trainer.ParseToken(token); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
	}
}

func (h *Handler) userIdentity(c *gin.Context) {
	token, ok := validHeader(c)
	if !ok {
		return
	}

	if err := h.services.User.ParseToken(token); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
	}
}
