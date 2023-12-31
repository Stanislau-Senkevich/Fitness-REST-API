package handler

import (
	"Fitness_REST_API/internal/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Get users full info
// @Security ApiKeyAuth
// @Tags admin
// @Description get full information about all users (not trainers)
// @ID get-users-full-info
// @Produce  json
// @Success 200 {object} usersInfoResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /admin/user [get]
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

// @Summary Get trainers full info
// @Security ApiKeyAuth
// @Tags admin
// @Description get full information about all trainers
// @ID get-trainers-full-info
// @Produce  json
// @Success 200 {object} usersInfoResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /admin/trainer [get]
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

// @Summary Get user full info
// @Security ApiKeyAuth
// @Tags admin
// @Description get full information about user by id
// @ID get-user-full-info
// @Produce  json
// @Success 200 {object} entity.UserInfo
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /admin/user/:id [get]
func (h *Handler) getUserFullInfoByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
		return
	}

	userInfo, err := h.services.GetUserFullInfoById(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, userInfo)
}

// @Summary Create user
// @Security ApiKeyAuth
// @Tags admin
// @Description creates user
// @ID create-user
// @Accept  json
// @Produce  json
// @Param input body entity.User true "user info"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /admin/user [post]
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
		"id": id,
	})
}

// @Summary Update user
// @Security ApiKeyAuth
// @Tags admin
// @Description updates user info
// @ID update-user
// @Accept  json
// @Produce  json
// @Param input body entity.UserUpdate true "update info"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /admin/user/:id [put]
func (h *Handler) updateUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
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

// @Summary Delete user
// @Security ApiKeyAuth
// @Tags admin
// @Description deletes user
// @ID delete-user
// @Accept  json
// @Produce  json
// @Param input body entity.UserUpdate true "update info"
// @Success 200 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /admin/user/:id [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
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
