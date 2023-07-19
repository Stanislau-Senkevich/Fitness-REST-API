package handler

import (
	"Fitness_REST_API/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type staffSignInInput struct {
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

type signUpResponse struct {
	Id int64 `json:"id"`
}

func (h *Handler) signUp(c *gin.Context) {
	var inputUser entity.User

	if err := c.ShouldBindJSON(&inputUser); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.services.User.SignUp(&inputUser)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, signUpResponse{
		Id: id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input userSignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	token, err := h.services.User.SignIn(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: token,
	})
}

func (h *Handler) adminSignIn(c *gin.Context) {
	var input staffSignInInput
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

func (h *Handler) trainerSignIn(c *gin.Context) {
	var input staffSignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	token, err := h.services.Trainer.SignIn(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: token,
	})
}
