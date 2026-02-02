package initializers

import (
	"github.com/gofiber/fiber/v2"
)

func NewFiberConfig() fiber.Config {
	cfg := fiber.Config{
		ServerHeader:  "EMob",
		AppName:       "EMob App v0.1-beta",
		CaseSensitive: true,
	}

	return cfg
}
