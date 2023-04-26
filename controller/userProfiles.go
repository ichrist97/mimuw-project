package handler

import (
	"fmt"
	"strconv"

	db "mimuw-project/database"
	model "mimuw-project/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserProfiles(c *fiber.Ctx) error {
	var cookie = c.Params("cookie")
	fmt.Println(cookie)
	var timeRangeStr = c.Query("time_range")

	// TODO check correct time range format
	if timeRangeStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Time range required")
	}
	// check and parse time range
	//var timeRangeSplit = strings.Split(timeRangeStr, "_")
	//var timeFormat = "2006-01-02T15:04:05.000Z" // exactly this format
	/*
		var timestampFrom, err0 = time.Parse(timeFormat, timeRangeSplit[0])
		var timestampEnd, err1 = time.Parse(timeFormat, timeRangeSplit[1])
		if err0 != nil || err1 != nil {
			fmt.Println("Failed parsing time range")
			return c.SendStatus(500)
		}
	*/

	//fmt.Println(timestampFrom)
	//fmt.Println(timestampEnd)

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
	coll := db.DB.Database("mimuw").Collection("user_tags")

	// TODO sort user tags in descending time order
	// TODO filter by time

	// get views
	viewsFilter := bson.D{
		{
			"$and",
			bson.A{
				bson.D{{"cookie", bson.D{{"$eq", cookie}}}},
				bson.D{{"action", bson.D{{"$eq", "VIEW"}}}},
			},
		},
	}

	viewsCursor, viewsErr := coll.Find(db.Ctx, viewsFilter)
	if viewsErr != nil {
		fmt.Println("Failed to get user profiles")
		fmt.Println("Errors", err.Error())
		return c.SendStatus(500)
	}

	// parse into struct
	var viewsResults []model.UserTagEvent
	if err = viewsCursor.All(db.Ctx, &viewsResults); err != nil {
		fmt.Println("Errors", err.Error())
		return c.SendStatus(500)
	}

	// get buys
	buysFilter := bson.D{
		{
			"$and",
			bson.A{
				bson.D{{"cookie", bson.D{{"$eq", cookie}}}},
				bson.D{{"action", bson.D{{"$eq", "BUY"}}}},
			},
		},
	}

	buysCursor, buysErr := coll.Find(db.Ctx, buysFilter)
	if buysErr != nil {
		fmt.Println("Failed to get user profiles")
		fmt.Println("Errors", err.Error())
		return c.SendStatus(500)
	}

	// parse into struct
	var buysResults []model.UserTagEvent
	if err = buysCursor.All(db.Ctx, &buysResults); err != nil {
		fmt.Println("Errors", err.Error())
		return c.SendStatus(500)
	}

	if err != nil {
		fmt.Println("Failed to get user profiles")
		fmt.Println("Errors", err.Error())
		return c.SendStatus(500)
	}

	// create response
	userProfile := model.UserProfile{Cookie: cookie, Views: viewsResults, Buys: buysResults}
	return c.JSON(userProfile)
}
