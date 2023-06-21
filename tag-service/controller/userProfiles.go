package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	db "tag-service/database"
	model "tag-service/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func logResponses(c *fiber.Ctx, res model.UserProfile) {
	// debug response from api
	body := new(model.UserProfile)
	c.BodyParser(&body)

	// actual generated response
	coll := db.DB.Database("mimuw").Collection("log_profiles")

	doc := map[string]interface{}{
		"true":      body,
		"generated": res,
	}

	_, err := coll.InsertOne(db.Ctx, doc)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetUserProfiles(c *fiber.Ctx, debug bool) error {
	var cookie = c.Params("cookie")
	var timeRangeStr = c.Query("time_range")

	if timeRangeStr == "" {
		fmt.Println("Bad request: time range required")
		return c.Status(fiber.StatusBadRequest).SendString("Time range required")
	}
	// check and parse time range
	var timeRangeSplit = strings.Split(timeRangeStr, "_")

	// cast to time object
	var timeFormat = "2006-01-02T15:04:05.000" // exactly this format
	var timestampFrom, err0 = time.Parse(timeFormat, timeRangeSplit[0])
	var timestampEnd, err1 = time.Parse(timeFormat, timeRangeSplit[1])
	if err0 != nil || err1 != nil {
		fmt.Println("Failed parsing time range")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid time range")
	}

	// parse limit
	var limitStr = c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && len(limitStr) > 0 {
		fmt.Println("Bad request: Limit must be integer")
		return c.Status(fiber.StatusBadRequest).SendString("Limit must be integer")
	}
	if limit == 0 {
		limit = 200 // default
	}

	// read user tags for cookie from database
	coll := db.DB.Database("mimuw").Collection("user_tags")

	// set limit and sort descending by timestampg
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{"time", -1}})

	// get views
	viewsFilter := bson.D{
		{
			"$and",
			bson.A{
				bson.D{{"cookie", bson.D{{"$eq", cookie}}}},
				bson.D{{"action", bson.D{{"$eq", "VIEW"}}}},
				bson.D{{"time", bson.D{{"$gte", primitive.NewDateTimeFromTime(timestampFrom)}}}},
				bson.D{{"time", bson.D{{"$lte", primitive.NewDateTimeFromTime(timestampEnd)}}}},
			},
		},
	}

	viewsCursor, viewsErr := coll.Find(db.Ctx, viewsFilter, opts)
	if viewsErr != nil {
		fmt.Println(err.Error())
		return c.SendStatus(500)
	}

	// parse into struct
	var viewsResults []model.UserTagEvent
	if err = viewsCursor.All(db.Ctx, &viewsResults); err != nil {
		fmt.Println(err.Error())
		return c.SendStatus(500)
	}
	// empty result
	if viewsResults == nil {
		viewsResults = []model.UserTagEvent{}
	}

	// get buys
	buysFilter := bson.D{
		{
			"$and",
			bson.A{
				bson.D{{"cookie", bson.D{{"$eq", cookie}}}},
				bson.D{{"action", bson.D{{"$eq", "BUY"}}}},
				bson.D{{"time", bson.D{{"$gte", primitive.NewDateTimeFromTime(timestampFrom)}}}},
				bson.D{{"time", bson.D{{"$lte", primitive.NewDateTimeFromTime(timestampEnd)}}}},
			},
		},
	}

	buysCursor, buysErr := coll.Find(db.Ctx, buysFilter, opts)
	if buysErr != nil {
		fmt.Println(err.Error())
		return c.SendStatus(500)
	}

	// parse into struct
	var buysResults []model.UserTagEvent
	if err = buysCursor.All(db.Ctx, &buysResults); err != nil {
		fmt.Println(err.Error())
		return c.SendStatus(500)
	}
	// empty result
	if buysResults == nil {
		buysResults = []model.UserTagEvent{}
	}

	if err != nil {
		fmt.Println(err.Error())
		return c.SendStatus(500)
	}

	// create response
	userProfile := model.UserProfile{Cookie: cookie, Views: viewsResults, Buys: buysResults}

	// log response and expected result
	if debug {
		logResponses(c, userProfile)
	}

	return c.JSON(userProfile)
}
