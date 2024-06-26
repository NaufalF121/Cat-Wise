package server

import (
	"Gateway/internal/database"
	"github.com/gofiber/fiber/v2"

	_ "Gateway/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "Gateway",
			AppName:      "Gateway",
		}),

		db: database.New(),
	}

	return server
}
