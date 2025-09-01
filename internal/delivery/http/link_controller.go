package http

import (
	"devshort-backend/internal/delivery/http/middleware"
	"devshort-backend/internal/model"
	"devshort-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type LinkController struct {
	UseCase *usecase.LinkUseCase
	Log     *logrus.Logger
}

func NewLinkController(useCase *usecase.LinkUseCase, log *logrus.Logger) *LinkController {
	return &LinkController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *LinkController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateLinkRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}
	request.UserId = auth.ID

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating link")
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(model.WebResponse[*model.LinkResponse]{Data: response})
}

func (c *LinkController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetLinkRequest{
		UserId: auth.ID,
		ID:     ctx.Params("linkId"),
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting link")
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(model.WebResponse[*model.LinkResponse]{Data: response})
}

func (c *LinkController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.ListLinkRequest{
		UserId: auth.ID,
	}

	responses, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list links")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.LinkResponse]{Data: responses})
}

func (c *LinkController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateLinkRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ID = ctx.Params("linkId")
	
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating link")
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(model.WebResponse[*model.LinkResponse]{Data: response})
}

func (c *LinkController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	linkId := ctx.Params("linkId")

	request := &model.DeleteLinkRequest{
		UserId: auth.ID,
		ID:     linkId,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("error deleting link")
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}