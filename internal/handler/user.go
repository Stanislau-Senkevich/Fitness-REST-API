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

// @Summary Get user info
// @Security ApiKeyAuth
// @Tags user
// @Description get information about yourself
// @ID get-user-info
// @Produce  json
// @Success 200 {object} entity.User
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user [get]
func (h *Handler) getUserInfo(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	user, err := h.services.GetUserInfoById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

// @Summary Get all workouts
// @Security ApiKeyAuth
// @Tags user
// @Description get information about your workouts
// @ID get-user-workouts
// @Produce  json
// @Success 200 {object} workoutsResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/workout [get]
func (h *Handler) getUserWorkouts(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	w, err := h.services.User.GetUserWorkouts(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, workoutsResponse{
		Workouts: w,
	})
}

// @Summary Get workout by id
// @Security ApiKeyAuth
// @Tags user
// @Description get information about workout using workout id
// @ID get-workout-user
// @Produce  json
// @Success 200 {object} entity.Workout
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/workout/:id [get]
func (h *Handler) getWorkoutByIdForUser(c *gin.Context) {
	userId, err := getId(c)
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

// @Summary Get all trainers
// @Security ApiKeyAuth
// @Tags user
// @Description get information about all trainers
// @ID get-trainers
// @Produce  json
// @Success 200 {object} usersResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/trainer [get]
func (h *Handler) getAllTrainers(c *gin.Context) {
	trainers, err := h.services.GetTrainers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, usersResponse{
		Users: trainers,
	})
}

// @Summary Get trainer
// @Security ApiKeyAuth
// @Tags user
// @Description get information about trainer using id
// @ID get-trainer-by-id
// @Produce  json
// @Success 200 {object} entity.User
// @Failure 400 {object} errorResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/trainer/:id [get]
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

// @Summary Get partnerships
// @Security ApiKeyAuth
// @Description get information about your partnerships
// @Tags user
// @ID get-partnership
// @Produce  json
// @Success 200 {object} partnershipsResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/partnership [get]
func (h *Handler) getPartnerships(c *gin.Context) {
	userId, err := getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	partnerships, err := h.services.GetUserPartnerships(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, partnershipsResponse{
		Partnerships: partnerships,
	})
}

// @Summary Create workout
// @Security ApiKeyAuth
// @Description creates workout
// @Tags user
// @ID create-workout-as-user
// @Accept  json
// @Produce  json
// @Param input body entity.Workout true "workout info"
// @Success 200 {object} workoutIdResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/workout [post]
func (h *Handler) createUserWorkout(c *gin.Context) {
	userId, err := getId(c)
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
	c.JSON(http.StatusOK, workoutIdResponse{
		WorkoutId: workoutId,
	})
}

// @Summary Update workout
// @Security ApiKeyAuth
// @Description updates workout
// @Tags user
// @ID update-workout-user
// @Accept  json
// @Produce  json
// @Param input body entity.UpdateWorkout true "update workout info"
// @Success 200 {object} workoutIdResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/workout/:id [put]
func (h *Handler) updateWorkoutForUser(c *gin.Context) {
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workoutId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
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
// @Tags user
// @ID delete-workout-user
// @Produce  json
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/workout/:id [delete]
func (h *Handler) deleteWorkoutForUser(c *gin.Context) {
	userId, err := getId(c)
	if err != nil || userId < 1 {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	workoutId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || workoutId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}
	err = h.services.DeleteWorkout(workoutId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Send request
// @Security ApiKeyAuth
// @Description sends request to trainer to become his client
// @Tags user
// @ID send-request
// @Produce  json
// @Success 200 {object} requestIdResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/partnership/trainer/:id [post]
func (h *Handler) sendRequestToTrainer(c *gin.Context) {
	trainerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}
	userId, err := getId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	requestId, err := h.services.SendRequestToTrainer(trainerId, userId)
	if err != nil {
		if requestId == -1 {
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, requestIdResponse{
		RequestId: requestId,
	})
}

// @Summary End partnership
// @Security ApiKeyAuth
// @Description ends partnership with trainer
// @Tags user
// @ID end-partnership-as-user
// @Produce  json
// @Success 200 {object} partnershipIdResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/partnership/trainer/:id [put]
func (h *Handler) endPartnershipWithTrainer(c *gin.Context) {
	trainerId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || trainerId < 1 {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}
	userId, err := getId(c)
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

	c.JSON(http.StatusOK, partnershipIdResponse{
		PartnershipId: pId,
	})
}

func getId(c *gin.Context) (int64, error) {
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
