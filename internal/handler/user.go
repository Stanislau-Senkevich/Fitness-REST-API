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
	userIdCtx = "userId"
)

func (h *Handler) getUserInfo(c *gin.Context) {
	id, err := h.getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	user, err := h.services.GetUserInfoById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) createUserWorkout(c *gin.Context) {
	userId, err := h.getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var input entity.Workout
	err = initUserWorkout(c, &input, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
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
	id, err := h.getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	w, err := h.services.User.GetUserWorkouts(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *Handler) getWorkoutById(c *gin.Context) {
	userId, err := h.getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workoutId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	workout, err := h.services.User.GetWorkoutById(workoutId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, workout)
}

func (h *Handler) updateWorkout(c *gin.Context) {
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workoutId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id parameter"))
		return
	}

	userId, err := h.getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	input, err := h.initUpdateWorkout(c, workoutId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.User.UpdateWorkout(workoutId, userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": workoutId,
	})

}

func (h *Handler) deleteWorkout(c *gin.Context) {
	userId, err := h.getId(c)
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
	c.Status(http.StatusOK)
}

func (h *Handler) getAllTrainers(c *gin.Context) {
	trainers, err := h.services.GetTrainers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, trainers)
}

func (h *Handler) getTrainerById(c *gin.Context) {
	trainerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}
	trainer, err := h.services.GetTrainerById(trainerId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("no trainer was found on provided id"))
		return
	}
	c.JSON(http.StatusOK, trainer)
}

func (h *Handler) getPartnerships(c *gin.Context) {
	userId, err := h.getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	partnerships, err := h.services.GetUserPartnerships(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, partnerships)
}

func (h *Handler) sendRequestToTrainer(c *gin.Context) {
	trainerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}
	userId, err := h.getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	pId, err := h.services.SendRequestToTrainer(trainerId, userId)
	if err != nil {
		if pId == -1 {
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"request_id": pId,
	})
}

func (h *Handler) endPartnershipWithTrainer(c *gin.Context) {
	trainerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}
	userId, err := h.getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	pId, err := h.services.EndPartnershipWithTrainer(trainerId, userId)
	if err != nil {
		if pId == -1 {
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"partnership_id": pId,
	})
}

func (h *Handler) getId(c *gin.Context) (int64, error) {
	id, ok := c.Get(userIdCtx)
	if !ok {
		id = c.Request.Context().Value(userIdCtx)
		idInt, ok := id.(int64)
		if !ok || idInt < 1 {
			err := fmt.Errorf("invalid id")
			return -1, err
		}
		return idInt, nil
	}

	idInt, ok := id.(int64)
	if !ok || idInt < 1 {
		err := fmt.Errorf("invalid id")
		return -1, err
	}
	return idInt, nil
}

func (h *Handler) initUpdateWorkout(c *gin.Context, workoutId, userId int64) (*entity.UpdateWorkout, error) {
	var input entity.UpdateWorkout
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return nil, err
	}

	workout, err := h.services.GetWorkoutById(workoutId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return nil, err
	}

	if input.Title == "" {
		input.Title = workout.Title
	}

	if input.Date.IsZero() {
		input.Date = workout.Date
	}
	return &input, nil
}

func initUserWorkout(c *gin.Context, input *entity.Workout, userId int64) error {
	if err := c.ShouldBindJSON(&input); err != nil {
		return err
	}

	if input.UserId == 0 {
		input.UserId = userId
	}

	if input.UserId != userId {
		return errors.New("user_id from token and user_id from workout must match")
	}
	return nil
}
