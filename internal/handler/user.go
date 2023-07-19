package handler

import (
	"Fitness_REST_API/internal/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

	var input entity.UserWorkout

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

	workoutId, err := h.services.User.CreateWorkout(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": workoutId,
	})
}

func (h *Handler) getAllWorkouts(c *gin.Context) {

}

func (h *Handler) getWorkoutByID(c *gin.Context) {

}

func (h *Handler) updateWorkout(c *gin.Context) {

}

func (h *Handler) deleteWorkout(c *gin.Context) {

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
