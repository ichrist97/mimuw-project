package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	db "mimuw-project/database"
	handler "mimuw-project/handler"
	model "mimuw-project/model"
	"reflect"
)

func main() {
	app := fiber.New()

	// Initialize default config
	app.Use(logger.New())

	// init cassandra connection
	CassandraSession := db.Session
	defer CassandraSession.Close()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Allezon!")
	})

	// use case 1 adding user tags
	app.Post("/user_tags", model.ValidateUserTagEvent, func(c *fiber.Ctx) error {
		// validate request body
		body := new(model.UserTagEvent)
		c.BodyParser(&body)
		return handler.AddUserTag(c, body)
	})

	// POST /user_profiles/{cookie}?time_range=<time_range>?limit=<limit>
	// use case 2 getting user profiles
	app.Post("/user_profiles/:cookie", func(c *fiber.Ctx) error {
		// read query params
		var timeRange = c.Query("time_range")
		var limit = c.Query("limit")

		fmt.Println(reflect.TypeOf(timeRange))
		fmt.Println(limit)

		return c.SendStatus(204)
	})

	// POST /aggregates?time_range=<time_from>_<time_to>&action=BUY&brand_id=Nike&aggregates=COUNT
	// use case 3 aggregated user actions / statistics
	app.Post("/aggregates", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	// Last middleware to match anything
	/*
		app.Use(func(c *fiber.Ctx) {
			c.SendStatus(404) // => 404 "Not Found"
		})
	*/

	app.Listen(":3000")
}
