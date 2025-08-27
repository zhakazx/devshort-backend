package usecase

import (
	"context"
	"devshort-backend/internal/entity"
	"devshort-backend/internal/gateway/messaging"
	"devshort-backend/internal/model"
	"devshort-backend/internal/model/converter"
	"devshort-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LinkUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	LinkRepository *repository.LinkRepository
	UserRepository *repository.UserRepository
	LinkProducer   *messaging.LinkProducer
}

func NewLinkUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	linkRepository *repository.LinkRepository, userRepository *repository.UserRepository, linkProducer *messaging.LinkProducer) *LinkUseCase {
	return &LinkUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		LinkRepository: linkRepository,
		UserRepository: userRepository,
		LinkProducer:   linkProducer,
	}
}

func (c *LinkUseCase) Create(ctx context.Context, request *model.CreateLinkRequest) (*model.LinkResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find user")
		return nil, fiber.ErrNotFound
	}

	link := &entity.Link{
		ID:       uuid.NewString(),
		UserId:   user.ID,
		Title:    request.Title,
		ShortUrl: request.ShortUrl,
		LongUrl:  request.LongUrl,
		IsActive: request.IsActive,
	}

	if err := c.LinkRepository.Create(tx, link); err != nil {
		c.Log.WithError(err).Error("failed to create link")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	if c.LinkProducer != nil {
		event := converter.LinkToEvent(link)
		if err := c.LinkProducer.Send(event); err != nil {
			c.Log.WithError(err).Error("failed to publish link created event")
			return nil, fiber.ErrInternalServerError
		}
		c.Log.Info("Published link created event")
	} else {
		c.Log.Info("Kafka producer is disabled, skipping link created event")
	}

	return converter.LinkToResponse(link), nil
}

func (c *LinkUseCase) Get(ctx context.Context, req *model.GetLinkRequest) (*model.LinkResponse, error) {
	link := new(entity.Link)
	if err := c.LinkRepository.FindByIdAndUserId(c.DB.WithContext(ctx), link, req.ID, req.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find link")
		return nil, fiber.ErrNotFound
	}

	return converter.LinkToResponse(link), nil
}

func (c *LinkUseCase) List(ctx context.Context, request *model.ListLinkRequest) ([]model.LinkResponse, error) {
	links, err := c.LinkRepository.FindAllByUserId(c.DB.WithContext(ctx), request.UserId)
	if err != nil {
		c.Log.WithError(err).Error("failed to find links by user id")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.LinkResponse, len(links))
	for i, link := range links {
		responses[i] = *converter.LinkToResponse(&link)
	}

	return responses, nil
}

func (c *LinkUseCase) Update(ctx context.Context, request *model.UpdateLinkRequest) (*model.LinkResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	link := new(entity.Link)
	if err := c.LinkRepository.FindById(tx, link, request.ID); err != nil {
		c.Log.WithError(err).Error("failed to find link by id")
		return nil, fiber.ErrNotFound
	}

	link.Title = request.Title
	link.ShortUrl = request.ShortUrl
	link.LongUrl = request.LongUrl
	link.IsActive = request.IsActive

	if err := c.LinkRepository.Update(tx, link); err != nil {
		c.Log.WithError(err).Error("failed to update link")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	if c.LinkProducer != nil {
		event := converter.LinkToEvent(link)
		if err := c.LinkProducer.Send(event); err != nil {
			c.Log.WithError(err).Error("failed to publish link updated event")
			return nil, fiber.ErrInternalServerError
		}
		c.Log.Info("Published link updated event")
	} else {
		c.Log.Info("Kafka producer is disabled, skipping link updated event")
	}

	return converter.LinkToResponse(link), nil
}

func (c *LinkUseCase) Delete(ctx context.Context, request *model.DeleteLinkRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	link := new(entity.Link)
	if err := c.LinkRepository.FindByIdAndUserId(tx, link, request.ID, request.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find link")
		return fiber.ErrNotFound
	}

	if err := c.LinkRepository.Delete(tx, link); err != nil {
		c.Log.WithError(err).Error("failed to delete link")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}