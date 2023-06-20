package main

import (
	"fmt"
	"os"
	controller "tag-service/controller"
	model "tag-service/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// activate debug mode
	DEBUG_MODE := os.Getenv("DEBUG")
	DEBUG := false // default
	if DEBUG_MODE == "1" {
		DEBUG = true
		fmt.Println("DEBUG MODE")
	}

	app := fiber.New()

	// Initialize default config
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Allezon!")
	})

	// use case 1 adding user tags
	app.Post("/user_tags", model.ValidateUserTagEvent, func(c *fiber.Ctx) error {
		// validate request body
		body := new(model.UserTagEvent)
		c.BodyParser(&body)
		return controller.AddUserTag(c, body)
	})

	// POST /user_profiles/{cookie}?time_range=<time_range>?limit=<limit>
	// use case 2 getting user profiles
	app.Post("/user_profiles/:cookie", func(c *fiber.Ctx) error {
		return controller.GetUserProfiles(c, DEBUG)
	})

	// POST /aggregates?time_range=<time_from>_<time_to>&action=BUY&brand_id=Nike&aggregates=COUNT
	// use case 3 aggregated user actions / statistics
	app.Post("/aggregates", func(c *fiber.Ctx) error {
		return controller.GetAggregate(c, DEBUG)
	})

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000" // default
	}

	app.Listen(fmt.Sprintf(":%s", PORT))
}