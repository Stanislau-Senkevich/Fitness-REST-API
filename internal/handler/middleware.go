package handler

import (
	"Fitness_REST_API/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func validHeader(c *gin.Context) (string, bool) {
	header := c.Request.Header.Get("Authorization")

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, ErrorEmptyAuthHeader)
		return "", false
	}

	parts := strings.Fields(header)

	if len(parts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, ErrorInvalidAuthHeader)
		return "", false
	}

	if parts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, ErrorInvalidAuthHeader)
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
		return
	}
	c.String(http.StatusOK, "ok")
}

func (h *Handler) trainerIdentity(c *gin.Context) {
	token, ok := validHeader(c)
	if !ok {
		return
	}

	id, role, err := h.services.User.ParseToken(token)

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	if role != entity.TrainerRole {
		newErrorResponse(c, http.StatusForbidden, ErrorForbidden)
		return
	}

	c.Set(userIdCtx, id)
}

func (h *Handler) userIdentity(c *gin.Context) {
	token, ok := validHeader(c)
	if !ok {
		return
	}

	id, _, err := h.services.User.ParseToken(token)

	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(userIdCtx, id)
}
