package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/MaximSvinin/tape/pkg/model"
)

type infoUsecase interface {
	Info() (*model.TapeInfo, error)
}

type info struct {
	usecase infoUsecase
}

// @Summary		get tape info
// @Description	get tape info
// @Tags			info
// @Accept			json
// @Produce		json
// @Success		200	{object}	model.TapeInfo
// @Failure		500
// @Router			/info [get]
func (i *info) GetInfo(ctx *fiber.Ctx) error {
	info, err := i.usecase.Info()
	if err != nil {
		return err
	}

	return ctx.JSON(info)
}
