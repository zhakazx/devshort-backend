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
	linkRepository := repository.NewLinkRepository(config.Log)

	// setup producer
	var userProducer *messaging.UserProducer
	var linkProducer *messaging.LinkProducer

	if config.Producer != nil {
		userProducer = messaging.NewUserProducer(config.Producer, config.Log)
		linkProducer = messaging.NewLinkProducer(config.Producer, config.Log)
	}

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, userProducer)
	linkUseCase := usecase.NewLinkUseCase(config.DB, config.Log, config.Validate, linkRepository, userRepository,linkProducer)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	linkController := http.NewLinkController(linkUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth()

	routeConfig := route.RouteConfig{
		App:               config.App,
		UserController:    userController,
		LinkController:    linkController,
		AuthMiddleware:    authMiddleware,
	}
	routeConfig.Setup()
}
