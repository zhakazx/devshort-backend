package config

import (
	"devshort-backend/internal/delivery/http"
	"devshort-backend/internal/delivery/http/middleware"
	"devshort-backend/internal/delivery/http/route"
	"devshort-backend/internal/gateway/messaging"
	"devshort-backend/internal/repository"
	"devshort-backend/internal/usecase"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Producer sarama.SyncProducer
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)

	// setup producer
	var userProducer *messaging.UserProducer

	if config.Producer != nil {
		userProducer = messaging.NewUserProducer(config.Producer, config.Log)
	}

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, userProducer)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

	routeConfig := route.RouteConfig{
		App:               config.App,
		UserController:    userController,
		AuthMiddleware:    authMiddleware,
	}
	routeConfig.Setup()
}
