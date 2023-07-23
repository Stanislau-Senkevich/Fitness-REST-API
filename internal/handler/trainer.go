package handler

import (
	"Fitness_REST_API/internal/entity"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) getAllTrainerUsers(c *gin.Context) {
	id, err := h.getId(c)
	if err != nil || id < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	users, err := h.services.User.GetTrainerUsers(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) getAllTrainerRequests(c *gin.Context) {
	id, err := h.getId(c)
	if err != nil || id < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	users, err := h.services.User.GetTrainerRequests(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) getTrainerUser(c *gin.Context) {
	trainerId, err := h.getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}
	user, err := h.services.GetTrainerUserById(trainerId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) getTrainerRequest(c *gin.Context) {
	trainerId, err := h.getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	requestId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || requestId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}
	user, err := h.services.GetTrainerRequestById(requestId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) initPartnershipWithUser(c *gin.Context) {
	trainerId, err := h.getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	pId, err := h.services.InitPartnershipWithUser(trainerId, userId)
	if err != nil {
		if pId == 0 {
			newErrorResponse(c, http.StatusInternalServerError, err)
		} else {
			newErrorResponse(c, http.StatusBadRequest, err)
		}
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"partnership_id": pId,
	})

}

func (h *Handler) endPartnershipWithUser(c *gin.Context) {
	trainerId, err := h.getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	pId, err := h.services.EndPartnershipWithUser(trainerId, userId)
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

func (h *Handler) acceptRequest(c *gin.Context) {
	trainerId, err := h.getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	requestId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || requestId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	pId, err := h.services.AcceptRequest(trainerId, requestId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"partnership_id": pId,
	})
}

func (h *Handler) denyRequest(c *gin.Context) {
	trainerId, err := h.getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	requestId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || requestId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	err = h.services.DenyRequest(trainerId, requestId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) createTrainerWorkout(c *gin.Context) {
	trainerId, err := h.getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var input entity.Workout
	err = initTrainerWorkout(c, &input, trainerId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	workoutId, err := h.services.User.CreateWorkoutAsTrainer(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"workout_id": workoutId,
	})
}

func (h *Handler) getTrainerWorkouts(c *gin.Context) {
	trainerId, err := h.getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	workouts, err := h.services.GetTrainerWorkouts(trainerId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, workouts)
}

func (h *Handler) getTrainerWorkoutsWithUser(c *gin.Context) {
	trainerId, err := h.getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

}

func initTrainerWorkout(c *gin.Context, input *entity.Workout, trainerId int64) error {
	if err := c.ShouldBindJSON(&input); err != nil {
		return err
	}

	if !input.TrainerId.Valid {
		input.TrainerId = sql.NullInt64{Int64: trainerId, Valid: true}
	}

	if input.TrainerId.Int64 != trainerId {
		return errors.New("trainer_id from token and trainer_id from workout must match")
	}

	if input.Date.IsZero() {
		input.Date = time.Now()
	}

	if input.UserId < 1 {
		return errors.New("invalid user_id")
	}
	return nil
}
