package handler

import (
	_ "Fitness_REST_API/docs"
	"Fitness_REST_API/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
		admin.GET("/user", h.getAllUsersFullInfo)
		admin.GET("/user/:id", h.getUserFullInfoByID)
		admin.POST("/user", h.createUser)
		admin.PUT("/user/:id", h.updateUser)
		admin.DELETE("/user/:id", h.deleteUser)

		admin.GET("/trainer", h.getTrainersInfo)
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
		trainer.GET("/workout/:id", h.getWorkoutByIdForTrainer)
		trainer.GET("/workout/user/:id", h.getTrainerWorkoutsWithUser)
		trainer.PUT("/workout/:id", h.updateWorkoutForUser)
		trainer.DELETE("/workout/:id", h.deleteWorkoutForTrainer)
	}
}

func (h *Handler) initUserRoutes(router *gin.Engine) {
	user := router.Group("/user", h.userIdentity)
	{
		user.GET("/", h.getUserInfo)

		user.GET("/workout", h.getUserWorkouts)
		user.GET("/workout/:id", h.getWorkoutByIdForUser)
		user.POST("/workout", h.createUserWorkout)
		user.PUT("/workout/:id", h.updateWorkoutForUser)
		user.DELETE("/workout/:id", h.deleteWorkoutForTrainer)

		user.GET("/trainer", h.getAllTrainers)
		user.GET("/trainer/:id", h.getTrainerById)

		user.GET("/partnership", h.getPartnerships)
		user.POST("/partnership/trainer/:id", h.sendRequestToTrainer)
		user.PUT("/partnership/trainer/:id", h.endPartnershipWithTrainer)
	}
}
