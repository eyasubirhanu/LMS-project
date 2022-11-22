package route

import (
	"test/api/handlers"
	"test/api/repository"
	"test/api/service"

	"test/api/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func NewInitializedServer(configuration config.Config) *gin.Engine {
	// Configuration
	gin.SetMode(gin.ReleaseMode)
	// router := gin.New()
	router := gin.Default()
	database := config.NewSQLite(configuration)

	// setup gin cors
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// MaxAge: 12 * time.Hour,
	}))

	// Setup Proxies (optional)
	// You can comment this section
	err := router.SetTrustedProxies([]string{configuration.Get("APP_URL")})
	if err != nil {
		panic(err)
	}
	userRepository := repository.NewUserRepository()

	// Email Verification Setup
	emailVerificationRepository := repository.NewEmailVerificationRepository()
	emailVerificationService := service.NewEmailService(&emailVerificationRepository, &userRepository, database)

	// User Setup
	userService := service.NewUserService(&userRepository, database, &emailVerificationRepository)
	userController := handlers.NewUserController(&userService, &emailVerificationService)

	// Routing
	userController.Route(router)

	return router
}
