package http

import "github.com/gofiber/fiber/v2"

type tmUsecase interface {
	Erase() error
	Eject() error
}

type tm struct {
	usecase tmUsecase
}

// @Summary		erase tape
// @Description	erase tape
// @Tags			tm
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		500
// @Router			/erase [delete]
func (tm *tm) Erase(ctx *fiber.Ctx) error {
	err := tm.usecase.Erase()
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// @Summary		eject tape
// @Description	eject tape
// @Tags			tm
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		500
// @Router			/eject [post]
func (tm *tm) Eject(ctx *fiber.Ctx) error {
	err := tm.usecase.Eject()
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
