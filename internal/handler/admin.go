package handler

import (
	"Fitness_REST_API/internal/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) getAllUsersFullInfo(c *gin.Context) {
	info := make([]*entity.UserInfo, 0)
	idSlice, err := h.services.GetUsersId(entity.UserRole)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	for _, id := range idSlice {
		userInfo, err := h.services.GetUserFullInfoById(id)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		info = append(info, userInfo)
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) getTrainersInfo(c *gin.Context) {
	info := make([]*entity.UserInfo, 0)
	idSlice, err := h.services.GetUsersId(entity.TrainerRole)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	for _, id := range idSlice {
		userInfo, err := h.services.GetUserFullInfoById(id)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		}
		info = append(info, userInfo)
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) getUserFullInfoByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	userInfo, err := h.services.GetUserFullInfoById(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, userInfo)
}

func (h *Handler) createUser(c *gin.Context) {
	var inputUser entity.User

	if err := c.ShouldBindJSON(&inputUser); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := checkRole(&inputUser); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.services.Admin.CreateUser(&inputUser)
	if err != nil {
		if id == -1 {
			newErrorResponse(c, http.StatusBadRequest, err)
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err)
		}
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": 1,
	})
}

func (h *Handler) updateUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	var update entity.UserUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.InitUpdateUser(userId, &update)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.Admin.UpdateUser(userId, &update)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) deleteUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	err = h.services.Admin.DeleteUser(userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}

func checkRole(user *entity.User) error {
	if user.Role == "" {
		user.Role = entity.UserRole
	}
	if user.Role != entity.UserRole && user.Role != entity.TrainerRole {
		return errors.New("invalid role")
	}
	return nil
}
