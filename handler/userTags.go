package handler

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	db "mimuw-project/database"
	model "mimuw-project/model"
)

func AddUserTag(c *fiber.Ctx, body *model.UserTagEvent) error {
	var errs []string
	var gocqlUuid gocql.UUID

	// have we created a user correctly
	var created bool = false

	// generate a unique UUID for this user
	gocqlUuid = gocql.TimeUUID()
	// write data to Cassandra
	if err := db.Session.Query(`
	INSERT INTO mimuwapi.userTagEvents (id, time, cookie, country, device, action, origin) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		gocqlUuid, body.Time, body.Cookie, body.Country, model.DeviceString(body.Device), model.ActionString(body.Action), body.Origin).Exec(); err != nil {
		errs = append(errs, err.Error())
	} else {
		created = true
	}

	// depending on whether we created the user, return the
	// resource ID in a JSON payload, or return our errors
	if created {
		fmt.Println("Created userTagEvent: ", gocqlUuid)
		return c.SendStatus(204)
	} else {
		fmt.Println("Failed to create user tag")
		fmt.Println("Errors", errs)
		return c.SendStatus(500)
	}
}
