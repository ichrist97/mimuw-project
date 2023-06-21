package main

import (
	"fmt"
	"os"
	controller "tag-service/controller"
	model "tag-service/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func debugMode() bool {
	// activate debug mode
	DEBUG_MODE := os.Getenv("DEBUG")
	DEBUG := false // default
	if DEBUG_MODE == "1" {
		DEBUG = true
		fmt.Println("DEBUG MODE")
	}
	return DEBUG
}

/*
func initKafkaProducer() (*kafka.Producer, string) {
	// read from env
	kafka_host := os.Getenv("KAFKA_HOST")
	if len(kafka_host) == 0 {
		kafka_host = "localhost:29092" // default
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if len(topic) == 0 {
		topic = "user_tags"
	}

	// get client id hostname
	client_hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Failed to get hostname: %s\n", err)
		os.Exit(1)
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafka_host,
		"client.id":         client_hostname,
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to kafka server")
	return p, topic
}
*/

func main() {
	DEBUG := debugMode()

	// init kafka producer
	//kafka_producer, topic := initKafkaProducer()
	/*
		var kafka_producer *kafka.Producer
		kafka_producer = nil
		topic := ""
	*/

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
		//return controller.AddUserTag(c, body, kafka_producer, topic)
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
