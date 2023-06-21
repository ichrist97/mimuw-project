package controller

import (
	"fmt"
	db "tag-service/database"
	model "tag-service/model"

	"github.com/gofiber/fiber/v2"
)

// func AddUserTag(c *fiber.Ctx, body *model.UserTagEvent, producer *kafka.Producer, kafka_topic string) error {
func AddUserTag(c *fiber.Ctx, body *model.UserTagEvent) error {
	// insert into mongodb
	coll := db.DB.Database("mimuw").Collection("user_tags")
	productInfo := model.Product{ProductId: body.ProductInfo.ProductId, BrandId: body.ProductInfo.BrandId, CategoryId: body.ProductInfo.CategoryId, Price: body.ProductInfo.Price}
	doc := model.UserTagEvent{Time: body.Time, Cookie: body.Cookie, Country: body.Country, Device: body.Device, Action: body.Action, Origin: body.Origin, ProductInfo: productInfo}
	_, err := coll.InsertOne(db.Ctx, doc)

	// sent to kafka topic
	//emitKafkaMessage(producer, kafka_topic, doc)

	// depending on whether we created the user, return the
	// resource ID in a JSON payload, or return our errors
	if err == nil {
		return c.SendStatus(204)
	} else {
		fmt.Println("Failed to create user tag")
		fmt.Println("Error: ", err.Error())
		return c.SendStatus(500)
	}
}

/*
func emitKafkaMessage(producer *kafka.Producer, topic string, doc model.UserTagEvent) {
	// conv doc to json string
	msg, err := json.Marshal(&doc)
	if err != nil {
		fmt.Printf("Failed to marshal user tag: %s\n", err)
		return
	}

	delivery_chan := make(chan kafka.Event, 10000)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg},
		delivery_chan,
	)
	if err != nil {
		fmt.Printf("Failed to send message to kafka: %s\n", err)
	}
}
*/
