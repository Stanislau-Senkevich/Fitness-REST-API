package handler

import (
	"Fitness_REST_API/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type adminSignInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userSignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signInResponse struct {
	Token string `json:"token"`
}

// @Summary Sign Up
// @Tags auth
// @Description creates user account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body entity.User true "account info"
// @Success 200 {object} idResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var inputUser entity.User

	if err := c.ShouldBindJSON(&inputUser); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.services.User.SignUp(&inputUser)
	if err != nil {
		if id == -1 {
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}
	c.JSON(http.StatusOK, idResponse{
		Id: id,
	})
}

// @Summary Sign In
// @Tags auth
// @Description sign-in
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param input body userSignInInput true "account info"
// @Success 200 {object} signInResponse
// @Failure 400,401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input userSignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	token, err := h.services.User.SignIn(input.Email, input.Password, entity.UserRole)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: token,
	})
}

// @Summary Sign In for admin
// @Tags auth
// @Description sign-in as admin
// @ID admin-sign-in
// @Accept  json
// @Produce  json
// @Param input body adminSignInInput true "account info"
// @Success 200 {object} signInResponse
// @Failure 400,401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /admin/auth/sign-in [post]
func (h *Handler) adminSignIn(c *gin.Context) {
	var input adminSignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	token, err := h.services.Admin.SignIn(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: token,
	})
}

// @Summary Sign In for trainer
// @Tags auth
// @Description sign-in as trainer
// @ID trainer-sign-in
// @Accept  json
// @Produce  json
// @Param input body userSignInInput true "account info"
// @Success 200 {object} signInResponse
// @Failure 400,401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trainer/auth/sign-in [post]
func (h *Handler) trainerSignIn(c *gin.Context) {
	var input userSignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	token, err := h.services.User.SignIn(input.Email, input.Password, entity.TrainerRole)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: token,
	})
}
