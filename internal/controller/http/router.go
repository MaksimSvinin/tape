package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"

	swaggo "github.com/MaximSvinin/tape/docs/swagger/tape"
	"github.com/MaximSvinin/tape/pkg/tape"
)

//go:generate go install github.com/swaggo/swag/cmd/swag@latest
//go:generate swag init -o ../../../docs/swagger/tape -g router.go --parseInternal --parseDependency  --instanceName tape

// @title			Tape api
// @version		1.0
// @description	Tape manage server.
func Router(tape tape.Tape) {
	app := fiber.New()

	{
		r := info{usecase: tape}
		app.Get("/info", r.GetInfo)
	}

	{
		r := tm{usecase: tape}

		app.Post("/eject", r.Eject)
		app.Delete("/erase", r.Erase)
	}

	{
		r := file{
			usecase: tape,
		}
		app.Post("/write", r.WriteFile)
		app.Get("/read", r.ReadFile)
	}

	app.Get("/swagger/*", swagger.New(swagger.Config{
		InstanceName: swaggo.SwaggerInfotape.InstanceName(),
	}))
	err := app.Listen("127.0.0.1:8080")
	if err != nil {
		log.Fatal().Err(err).Msg("error run server")
	}
}
