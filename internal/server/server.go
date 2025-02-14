package server

import (
	"github.com/gofiber/fiber/v2"

	"gapi-platform/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "gapi-platform",
			AppName:      "gapi-platform",
		}),

		db: database.New(),
	}

	return server
}
