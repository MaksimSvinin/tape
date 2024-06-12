package http

import (
	"errors"
	"io"
	"os"

	"github.com/MaximSvinin/tape/internal/entity"
	"github.com/MaximSvinin/tape/pkg/model"
	"github.com/gofiber/fiber/v2"
)

type fileUsecase interface {
	Write(file io.Reader) (*model.FileWriteInfo, error)
	Read(fileNumbers []int, patch string) error
}

type file struct {
	usecase fileUsecase
}

// @Summary		write file
// @Description	write file
// @Tags			file
// @Accept			json
// @Param			filePath	query	string	true	"путь до файла"
// @Produce		json
// @Success		200	{object}	model.FileWriteInfo
// @Failure		500
// @Router			/write [post]
func (f *file) WriteFile(ctx *fiber.Ctx) error {
	filePath := ctx.Query("filePath")
	if filePath == "" {
		err := ctx.SendStatus(fiber.StatusBadRequest)
		if err != nil {
			return err
		}
		return errors.New("no file path")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	resp, err := f.usecase.Write(file)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

// @Summary		read file
// @Description	read file
// @Tags			file
// @Accept			json
// @Param			numbers	query	[]int	false	"file numbers"
// @Param			path	query	string	true	"path to extract dir"
// @Success		200
// @Failure		500
// @Router			/read [get]
func (f *file) ReadFile(ctx *fiber.Ctx) error {
	query := new(entity.ReadFileQuery)
	err := ctx.QueryParser(query)
	if err != nil {
		return err
	}

	if query.Path == "" {
		err = ctx.SendStatus(fiber.StatusBadRequest)
		if err != nil {
			return err
		}
		return errors.New("no path to extract dir")
	}

	err = f.usecase.Read(query.Numbers, query.Path)
	if err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}
