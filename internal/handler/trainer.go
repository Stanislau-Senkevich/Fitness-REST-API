package handler

import (
	"Fitness_REST_API/internal/entity"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Get clients
// @Security ApiKeyAuth
// @Tags trainer
// @Description get information about trainer clients
// @ID get-trainer-clients
// @Produce  json
// @Success 200 {object} usersResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/user [get]
func (h *Handler) getTrainerUsers(c *gin.Context) {
	id, err := getId(c)
	if err != nil || id < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	users, err := h.services.User.GetTrainerUsers(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, usersResponse{
		Users: users,
	})
}

// @Summary Get requests
// @Security ApiKeyAuth
// @Tags trainer
// @Description get information about users which send request to trainer
// @ID get-trainer-requests
// @Produce  json
// @Success 200 {object} usersResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/request [get]
func (h *Handler) getTrainerRequests(c *gin.Context) {
	id, err := getId(c)
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

// @Summary Get user by id
// @Security ApiKeyAuth
// @Tags trainer
// @Description get information about trainer client
// @ID get-trainer-user-by-id
// @Produce  json
// @Success 200 {object} entity.User
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/user/:id [get]
func (h *Handler) getTrainerUserById(c *gin.Context) {
	trainerId, err := getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
		return
	}
	user, err := h.services.GetTrainerUserById(trainerId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Get request by id
// @Security ApiKeyAuth
// @Tags trainer
// @Description get information about trainer request by id
// @ID get-trainer-request
// @Produce  json
// @Success 200 {object} entity.Request
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/request/:id [get]
func (h *Handler) getTrainerRequestById(c *gin.Context) {
	trainerId, err := getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	requestId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || requestId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
		return
	}
	request, err := h.services.GetTrainerRequestById(trainerId, requestId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, request)
}

// @Summary Get workouts
// @Security ApiKeyAuth
// @Tags trainer
// @Description get information about trainer workouts
// @ID get-trainer-workouts
// @Produce  json
// @Success 200 {object} workoutsResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/workout [get]
func (h *Handler) getTrainerWorkouts(c *gin.Context) {
	trainerId, err := getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	workouts, err := h.services.GetTrainerWorkouts(trainerId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, workoutsResponse{
		Workouts: workouts,
	})
}

// @Summary Get workouts with user
// @Security ApiKeyAuth
// @Tags trainer
// @Description get information about trainer workouts with user
// @ID get-trainer-workouts-user
// @Produce  json
// @Success 200 {object} workoutsResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/workout/user/:id [get]
func (h *Handler) getTrainerWorkoutsWithUser(c *gin.Context) {
	trainerId, err := getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
		return
	}

	workouts, err := h.services.GetTrainerWorkoutsWithUser(trainerId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, workoutsResponse{
		Workouts: workouts,
	})
}

// @Summary Create workout
// @Security ApiKeyAuth
// @Tags trainer
// @Description creates workout with user
// @ID create-workout-as-trainer
// @Accept  json
// @Produce  json
// @Param input body entity.Workout true "workout info"
// @Success 200 {object} workoutIdResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/workout [post]
func (h *Handler) createTrainerWorkout(c *gin.Context) {
	trainerId, err := getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var input entity.Workout
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = formatTrainerWorkout(&input, trainerId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	workoutId, err := h.services.User.CreateWorkoutAsTrainer(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, workoutIdResponse{
		WorkoutId: workoutId,
	})
}

// @Summary Init partnership
// @Security ApiKeyAuth
// @Tags trainer
// @Description starts partnership with user if possible
// @ID init-partnership-with-user
// @Produce  json
// @Success 200 {object} partnershipIdResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/user/:id [post]
func (h *Handler) initPartnershipWithUser(c *gin.Context) {
	trainerId, err := getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
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

	c.JSON(http.StatusOK, partnershipIdResponse{
		PartnershipId: pId,
	})
}

// @Summary End partnership
// @Security ApiKeyAuth
// @Tags trainer
// @Description ends partnership with user if possible
// @ID end-partnership-with-user
// @Produce  json
// @Success 200 {object} partnershipIdResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/user/:id [put]
func (h *Handler) endPartnershipWithUser(c *gin.Context) {
	trainerId, err := getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
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

	c.JSON(http.StatusOK, partnershipIdResponse{
		PartnershipId: pId,
	})
}

// @Summary Accept request
// @Security ApiKeyAuth
// @Tags trainer
// @Description accepts request from user by provided request id if possible
// @ID accept-request
// @Produce  json
// @Success 200 {object} partnershipIdResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/request/:id [put]
func (h *Handler) acceptRequest(c *gin.Context) {
	trainerId, err := getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	requestId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || requestId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
		return
	}

	pId, err := h.services.AcceptRequest(trainerId, requestId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, partnershipIdResponse{
		PartnershipId: pId,
	})
}

// @Summary Deny request
// @Security ApiKeyAuth
// @Tags trainer
// @Description deny and delete request from user by provided request id
// @ID deny-request
// @Produce  json
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/request/:id [delete]
func (h *Handler) denyRequest(c *gin.Context) {
	trainerId, err := getId(c)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	requestId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || requestId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
		return
	}

	err = h.services.DenyRequest(trainerId, requestId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Get workout by id
// @Security ApiKeyAuth
// @Tags user, trainer
// @Description get information about workout using workout id
// @ID get-workout-trainer
// @Produce  json
// @Success 200 {object} entity.Workout
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/workout/:id [get]
func (h *Handler) getWorkoutByIdForTrainer(c *gin.Context) {
	userId, err := getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workoutId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
		return
	}

	workout, err := h.services.User.GetWorkoutById(workoutId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, workout)
}

// @Summary Update workout
// @Security ApiKeyAuth
// @Description updates workout
// @Tags trainer
// @ID update-workout-trainer
// @Accept  json
// @Produce  json
// @Param input body entity.UpdateWorkout true "update workout info"
// @Success 200 {object} workoutIdResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/workout/:id [put]
func (h *Handler) updateWorkoutForTrainer(c *gin.Context) {
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workoutId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
		return
	}

	userId, err := getId(c)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var input entity.UpdateWorkout
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.FormatUpdateWorkout(&input, workoutId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = h.services.User.UpdateWorkout(workoutId, userId, &input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, workoutIdResponse{
		WorkoutId: workoutId,
	})
}

// @Summary Delete workout
// @Security ApiKeyAuth
// @Description deletes workout
// @Tags trainer
// @ID delete-workout-trainer
// @Produce  json
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/workout/:id [delete]
func (h *Handler) deleteWorkoutForTrainer(c *gin.Context) {
	userId, err := getId(c)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workoutId < 1 {
		newErrorResponse(c, http.StatusBadRequest, ErrorInvalidIdParameter)
		return
	}
	err = h.services.DeleteWorkout(workoutId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}

func formatTrainerWorkout(input *entity.Workout, trainerId int64) error {
	if !input.TrainerId.Valid {
		input.TrainerId = sql.NullInt64{Int64: trainerId, Valid: true}
	}

	if input.TrainerId.Int64 != trainerId {
		return errors.New("trainer_id from token and trainer_id from workout must match")
	}

	if input.UserId < 1 {
		return errors.New("invalid user_id")
	}
	return nil
}
