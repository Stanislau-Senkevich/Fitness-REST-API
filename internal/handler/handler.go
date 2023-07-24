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
		admin.GET("/trainer/:id", h.getTrainerById)
		admin.PUT("/trainer/:id", h.updateTrainer)
		admin.DELETE("/trainer/:id", h.deleteTrainer)
		admin.POST("/trainer", h.createTrainer)
	}
}

func (h *Handler) initTrainerRoutes(router *gin.Engine) {
	trainer := router.Group("/trainer", h.trainerIdentity)
	{
		trainer.GET("/user", h.getTrainerUsers)
		trainer.GET("/user/:id", h.getTrainerUserById)
		trainer.POST("/user/:id", h.initPartnershipWithUser)
		trainer.PUT("/user/:id", h.endPartnershipWithUser)

		trainer.GET("/request", h.getTrainerRequests)
		trainer.GET("/request/:id", h.getTrainerRequestById)
		trainer.PUT("/request/:id", h.acceptRequest)
		trainer.DELETE("/request/:id", h.denyRequest)

		trainer.POST("/workout", h.createTrainerWorkout)
		trainer.GET("/workout", h.getTrainerWorkouts)
		trainer.GET("/workout/:id", h.getWorkoutById)
		trainer.GET("/workout/user/:id", h.getTrainerWorkoutsWithUser)
		trainer.PUT("/workout/:id", h.updateWorkout)
		trainer.DELETE("/workout/:id", h.deleteWorkout)
	}
}

func (h *Handler) initUserRoutes(router *gin.Engine) {
	user := router.Group("/user", h.userIdentity)
	{
		user.GET("/", h.getUserInfo)

		user.GET("/workout", h.getUserWorkouts)
		user.GET("/workout/:id", h.getWorkoutById)
		user.POST("/workout", h.createUserWorkout)
		user.PUT("/workout/:id", h.updateWorkout)
		user.DELETE("/workout/:id", h.deleteWorkout)

		user.GET("/trainer", h.getAllTrainers)
		user.GET("/trainer/:id", h.getTrainerById)

		user.GET("/partnership", h.getPartnerships)
		user.POST("/partnership/trainer/:id", h.sendRequestToTrainer)
		user.PUT("/partnership/trainer/:id", h.endPartnershipWithTrainer)
	}
}
