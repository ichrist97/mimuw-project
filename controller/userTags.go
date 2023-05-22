package controller

import (
	"fmt"
	db "mimuw-project/database"
	model "mimuw-project/model"

	"github.com/gofiber/fiber/v2"
)

func AddUserTag(c *fiber.Ctx, body *model.UserTagEvent) error {
	// insert into mongodb
	coll := db.DB.Database("mimuw").Collection("user_tags")
	doc := model.UserTagEvent{Time: body.Time, Cookie: body.Cookie, Country: body.Country, Device: body.Device, Action: body.Action, Origin: body.Origin}
	result, err := coll.InsertOne(db.Ctx, doc)

	// depending on whether we created the user, return the
	// resource ID in a JSON payload, or return our errors
	if err == nil {
		fmt.Println("Created userTagEvent: ", result.InsertedID)
		return c.SendStatus(204)
	} else {
		fmt.Println("Failed to create user tag")
		fmt.Println("Error: ", err.Error())
		return c.SendStatus(500)
	}
}
