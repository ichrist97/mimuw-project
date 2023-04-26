package handler

/*

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	//"github.com/gocql/gocql"
	db "mimuw-project/database"
	model "mimuw-project/model"

	"github.com/gofiber/fiber/v2"
)

func GetUserProfiles(c *fiber.Ctx) error {
	var cookie = c.Params("cookie")
	var timeRangeStr = c.Query("time_range")

	// TODO check correct time range format
	if timeRangeStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Time range required")
	}
	// check and parse time range
	var timeRangeSplit = strings.Split(timeRangeStr, "_")
	var timeFormat = "2006-01-02T15:04:05.000Z" // exactly this format
	var timestampFrom, err0 = time.Parse(timeFormat, timeRangeSplit[0])
	var timestampEnd, err1 = time.Parse(timeFormat, timeRangeSplit[1])
	if err0 != nil || err1 != nil {
		fmt.Println("Failed parsing time range")
		return c.SendStatus(500)
	}

	fmt.Println(timestampFrom)
	fmt.Println(timestampEnd)

	// parse limit
	var limitStr = c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && len(limitStr) > 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Limit must be integer")
	}
	if limit == 0 {
		limit = 200 // default
	}

	// read user tags for cookie from database
	dbErr := db.Session.Query(
		`SELECT * FROM mimuwapi.userTagEvents WHERE cookie = ? LIMIT ?`, cookie, limit).Exec()
	if err != nil {
		fmt.Println("Failed to get user profiles")
		fmt.Println("Errors", dbErr.Error())
		return c.SendStatus(500)
	}

	// hardcoded response
	u := model.UserProfile{Cookie: "cookie", Views: []model.UserTagEvent{}, Buys: []model.UserTagEvent{}}
	return c.JSON(u)
}

*/
