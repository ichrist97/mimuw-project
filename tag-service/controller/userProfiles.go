package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
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

func validateProfileReq(c *fiber.Ctx) (*string, *time.Time, *time.Time, *int, error) {
	var cookie = c.Params("cookie")
	var timeRangeStr = c.Query("time_range")

	if timeRangeStr == "" {
		fmt.Println("Bad request: time range required")
		return nil, nil, nil, nil, errors.New("time range required")
	}
	// check and parse time range
	var timeRangeSplit = strings.Split(timeRangeStr, "_")

	// cast to time object
	var timeFormat = "2006-01-02T15:04:05.000" // exactly this format
	var timestampFrom, err0 = time.Parse(timeFormat, timeRangeSplit[0])
	var timestampEnd, err1 = time.Parse(timeFormat, timeRangeSplit[1])
	if err0 != nil || err1 != nil {
		fmt.Println("Failed parsing time range")
		return nil, nil, nil, nil, err0
	}

	// parse limit
	var limitStr = c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil && len(limitStr) > 0 {
		fmt.Println("Bad request: Limit must be integer")
		return nil, nil, nil, nil, err
	}
	if limit == 0 {
		limit = 200 // default
	}

	// success
	return &cookie, &timestampFrom, &timestampEnd, &limit, nil
}

func profileWorker(id int, wg *sync.WaitGroup, filter *bson.D, limit *int, action string, resultChan chan<- *model.UserProfileWorkerResult) {
	defer wg.Done() // Signal the WaitGroup that this goroutine is done
	// Perform some work
	res, err := queryProfile(filter, limit)
	// Send the result through the channel
	chanRes := model.UserProfileWorkerResult{Results: res, Err: err, Action: action}
	resultChan <- &chanRes
}

func queryProfile(filter *bson.D, limit *int) (*[]model.UserTagEvent, error) {
	// read user tags for cookie from database
	coll := db.DB.Database("mimuw").Collection("user_tags")

	// set limit and sort descending by timestampg
	opts := options.Find().SetLimit(int64(*limit)).SetSort(bson.D{{"time", -1}})

	cursor, err := coll.Find(db.Ctx, *filter, opts)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// parse into struct
	var results []model.UserTagEvent
	if err = cursor.All(db.Ctx, &results); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	// empty result
	if results == nil {
		results = []model.UserTagEvent{}
	}
	return &results, nil
}

func GetUserProfiles(c *fiber.Ctx, debug bool) error {
	cookie, timestampFrom, timestampEnd, limit, err := validateProfileReq(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(2) // wait for two goroutines

	// Create a channel to receive the results
	resultChan := make(chan *model.UserProfileWorkerResult)

	// get views
	viewsFilter := bson.D{
		{
			"$and",
			bson.A{
				bson.D{{"cookie", bson.D{{"$eq", cookie}}}},
				bson.D{{"action", bson.D{{"$eq", "VIEW"}}}},
				bson.D{{"time", bson.D{{"$gte", primitive.NewDateTimeFromTime(*timestampFrom)}}}},
				bson.D{{"time", bson.D{{"$lte", primitive.NewDateTimeFromTime(*timestampEnd)}}}},
			},
		},
	}
	// start goroutine for views results
	go profileWorker(1, &wg, &viewsFilter, limit, "VIEW", resultChan)

	// get buys
	buysFilter := bson.D{
		{
			"$and",
			bson.A{
				bson.D{{"cookie", bson.D{{"$eq", cookie}}}},
				bson.D{{"action", bson.D{{"$eq", "BUY"}}}},
				bson.D{{"time", bson.D{{"$gte", primitive.NewDateTimeFromTime(*timestampFrom)}}}},
				bson.D{{"time", bson.D{{"$lte", primitive.NewDateTimeFromTime(*timestampEnd)}}}},
			},
		},
	}
	// start goroutine for buys results
	go profileWorker(2, &wg, &buysFilter, limit, "BUY", resultChan)

	// wait for goroutines to finish and close channels
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// collect results
	var viewsResults []model.UserTagEvent
	var buysResults []model.UserTagEvent
	for result := range resultChan {
		if result.Err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Err.Error())
		}
		if result.Action == "VIEW" {
			viewsResults = *result.Results
		} else if result.Action == "BUY" {
			buysResults = *result.Results
		}
	}

	// create response
	userProfile := model.UserProfile{Cookie: *cookie, Views: viewsResults, Buys: buysResults}

	// log response and expected result
	if debug {
		logResponses(c, userProfile)
	}

	return c.JSON(userProfile)
}
