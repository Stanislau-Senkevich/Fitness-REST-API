package handler

import (
	"Fitness_REST_API/internal/entity"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	userCtx = "userId"
)

func (h *Handler) getUserInfo(c *gin.Context) {
	id, err := h.getId(c, userCtx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	user, err := h.services.GetUser(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) createWorkout(c *gin.Context) {
	id, err := h.getId(c, userCtx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var input entity.Workout

	if err = c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if input.UserId == 0 {
		input.UserId = id
	}

	if input.UserId != id {
		newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("user_id from token and user_id from workout must match"))
		return
	}

	workoutId, err := h.services.User.CreateWorkoutAsUser(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": workoutId,
	})
}

func (h *Handler) getUserWorkouts(c *gin.Context) {
	id, err := h.getId(c, userCtx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	w, err := h.services.User.GetAllUserWorkouts(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *Handler) getWorkoutByID(c *gin.Context) {
	userId, err := h.getId(c, userCtx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	workout, err := h.services.User.GetWorkoutById(workoutId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, workout)
}

func (h *Handler) updateWorkout(c *gin.Context) {
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id parameter"))
		return
	}

	userId, err := h.getId(c, userCtx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var input entity.UpdateWorkout
	if err = c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.User.UpdateWorkout(workoutId, userId, &input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": workoutId,
	})

}

func (h *Handler) deleteWorkout(c *gin.Context) {
	userId, err := h.getId(c, userCtx)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	err = h.services.DeleteWorkout(workoutId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) getAllTrainers(c *gin.Context) {

}

func (h *Handler) getTrainerByID(c *gin.Context) {

}

func (h *Handler) sendRequestToTrainer(c *gin.Context) {

}

func (h *Handler) deletePartnershipWithTrainer(c *gin.Context) {

}

func (h *Handler) getId(c *gin.Context, key string) (int64, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		err := fmt.Errorf("user was not found")
		return -1, err
	}

	idInt, ok := id.(int64)
	if !ok {
		err := fmt.Errorf("invalid id type")
		return -1, err
	}
	return idInt, nil
}
