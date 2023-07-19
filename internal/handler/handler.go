package handler

import (
	"Fitness_REST_API/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	h.initAuthRoutes(router)
	h.initAdminRoutes(router)
	h.initTrainerRoutes(router)
	h.initUserRoutes(router)
	return router
}

func (h *Handler) initAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/admin/sign-in", h.adminSignIn)
		auth.POST("/trainer/sign-in", h.trainerSignIn)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
	}

}

func (h *Handler) initAdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin", h.adminIdentity)
	{
		admin.GET("/user", h.getAllUsers)
		admin.GET("/user/:id", h.getUserByID)
		admin.PUT("/user/:id", h.updateUser)
		admin.DELETE("/user/:id", h.deleteUser)
		admin.POST("/user", h.createUser)

		admin.GET("/trainer", h.getAllTrainers)
		admin.GET("/trainer/:id", h.getTrainerByID)
		admin.PUT("/trainer/:id", h.updateTrainer)
		admin.DELETE("/trainer/:id", h.deleteTrainer)
		admin.POST("/trainer", h.createTrainer)
	}
}

func (h *Handler) initTrainerRoutes(router *gin.Engine) {
	trainer := router.Group("/trainer", h.trainerIdentity)
	{
		trainer.GET("/user", h.getAllTrainerUsers)
		trainer.GET("/user/:id", h.getTrainerUser)
		trainer.POST("/user/:id", h.addUserToTrainerList)
		trainer.DELETE("/user/:id", h.deleteUserFromTrainerList)

		trainer.GET("/request", h.getAllTrainerRequests)
		trainer.GET("/request/:id", h.getTrainerRequest)
		trainer.POST("/request/:id", h.acceptRequest)
		trainer.DELETE("/request/:id", h.denyRequest)

		trainer.POST("/workout", h.createTrainerWorkout)
		trainer.GET("/workout", h.getTrainerWorkouts)
		trainer.GET("/workout/:id", h.getTrainerWorkoutByID)
		trainer.GET("/workout/user/:id", h.getTrainerWorkoutsWithUser)
		trainer.PUT("/workout/:id", h.updateTrainerWorkout)
		trainer.DELETE("/workout/:id", h.deleteTrainerWorkout)
	}
}

func (h *Handler) initUserRoutes(router *gin.Engine) {
	user := router.Group("/user", h.userIdentity)
	{
		user.GET("/", h.getUserInfo)

		user.GET("/workout", h.getUserWorkouts)
		user.GET("/workout/:id", h.getWorkoutByID)
		user.POST("/workout", h.createWorkout)
		user.PUT("/workout/:id", h.updateWorkout)
		user.DELETE("/workout/:id", h.deleteWorkout)

		user.GET("/trainer", h.getAllTrainers)
		user.GET("/trainer/:id", h.getTrainerByID)
		user.POST("/request/trainer/:id", h.sendRequestToTrainer)
		user.DELETE("/trainer/:id", h.deletePartnershipWithTrainer)
	}
}
