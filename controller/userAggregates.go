package controller

import (
	"fmt"
	"math"
	db "mimuw-project/database"
	model "mimuw-project/model"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func deleteOldEntries(timestampFrom time.Time) (*mongo.DeleteResult, error) {
	// delete entries older than 24 hours in logical time
	coll := db.DB.Database("mimuw").Collection("user_tags")
	filter := bson.D{{"time", bson.D{{"lt", timestampFrom.Add(-(time.Hour * 24))}}}}
	results, err := coll.DeleteMany(db.Ctx, filter)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func createTimeBoundaries(timestampFrom time.Time, timestampEnd time.Time) bson.A {
	timeDiffMin := int(math.Ceil(timestampEnd.Sub(timestampFrom).Minutes()))
	var boundaries bson.A
	boundaries = append(boundaries, timestampFrom.Format(time.RFC3339))

	var t = timestampFrom
	for i := 0; i < timeDiffMin-1; i++ {
		// one minute to each new bucket
		t = t.Add(time.Minute)
		boundaries = append(boundaries, t.Format(time.RFC3339))
	}
	boundaries = append(boundaries, timestampEnd.Format(time.RFC3339))
	return boundaries
}

func GetAggregate(c *fiber.Ctx, debug bool) error {
	// parse queries
	query := new(model.AggregateRequest)
	if err := c.QueryParser(query); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).SendString("Invalid query")
	}

	// check time range query
	var timeRangeStr = query.TimeRange
	if timeRangeStr == "" {
		fmt.Println("Bad request: time range required")
		return c.Status(fiber.StatusBadRequest).SendString("Time range required")
	}
	// check and parse time range
	var timeRangeSplit = strings.Split(timeRangeStr, "_")

	// cast to time object
	var timeFormat = "2006-01-02T15:04:05.000Z" // exactly this format
	var timestampFrom, err0 = time.Parse(timeFormat, timeRangeSplit[0])
	var timestampEnd, err1 = time.Parse(timeFormat, timeRangeSplit[1])
	fmt.Println(timestampFrom, timestampEnd)
	if err0 != nil || err1 != nil {
		fmt.Println("Failed parsing time range")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid time range")
	}

	// check action query
	var action = query.Action
	var validActions = []string{"BUY", "VIEW"}
	if action == "" || !contains(validActions, action) {
		fmt.Println("Bad request: valid action required")
		return c.Status(fiber.StatusBadRequest).SendString("Valid action required")
	}

	// check aggregates query
	var validAggregates = []string{"COUNT", "SUM_PRICE"}
	var useCount bool
	var useSum bool
	if len(query.Aggregates) == 0 {
		fmt.Println("Invalid aggregates")
		return c.Status(fiber.StatusBadRequest).SendString("At least one aggregation required")
	}
	for _, aggr := range query.Aggregates {
		if !contains(validAggregates, aggr) {
			fmt.Println("Bad request: invalid aggregates")
			return c.Status(fiber.StatusBadRequest).SendString("Invalid aggregates")
		}
		if aggr == "COUNT" {
			useCount = true
		} else if aggr == "SUM_PRICE" {
			useSum = true
		}
	}

	// define 1min boundary timestamps for bucket aggregation
	boundaries := createTimeBoundaries(timestampFrom, timestampEnd)

	// TODO write mongo query
	// group into 1 minute buckets
	// result data structure as table columns
	matchCriteria := bson.M{
		"action": action,
	}
	// optional parameters
	// BUG TODO filtering only by origin does not work
	if query.Origin != "" {
		matchCriteria["origin"] = query.Origin
	}
	if query.BrandId != "" {
		matchCriteria["product_info.brand_id"] = query.BrandId
	}
	if query.Origin != "" {
		matchCriteria["product_info.category_id"] = query.CategoryId
	}

	outputStage := bson.M{
		"tags": bson.M{"$push": bson.M{"time": "$time", "product_price": "$product_info.price", "action": "$action"}},
	}
	// add aggregations to mongo stage
	if useCount {
		outputStage["count"] = bson.M{"$sum": 1}
	}

	pipe := mongo.Pipeline{
		{{
			"$match", matchCriteria,
		}},
		{{"$bucket", bson.M{
			"groupBy":    "$time",
			"boundaries": boundaries,
			"default":    "Other",
			"output":     outputStage,
		}}},
	}

	// pass the pipeline to the Aggregate() method
	coll := db.DB.Database("mimuw").Collection("user_tags")
	cursor, err := coll.Aggregate(db.Ctx, pipe)
	if err != nil {
		fmt.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// TODO handle empty results as empty array instead of nil
	// decode the results
	var results []model.AggregateQuery
	if err = cursor.All(db.Ctx, &results); err != nil {
		fmt.Println(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// parse results into columns wise structure
	tableResults := transformToTable(results, useCount, useSum, action, query.Origin, query.BrandId, query.CategoryId)

	// TODO discard old data after query which is older than 24h logical
	/*
		_, err := deleteOldEntries(timestampFrom)
		if err != {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	*/

	// TODO Think about the number of possible aggregates,
	// how to calculate them and possibly store globally for later querying.

	// log response and expected result
	if debug {
		logAggrResponses(c, &tableResults)
	}

	return c.JSON(tableResults)
}

func logAggrResponses(c *fiber.Ctx, res *model.AggregateResult) {
	// debug response from api
	body := new(model.AggregateResult)
	c.BodyParser(&body)

	// actual generated response
	coll := db.DB.Database("mimuw").Collection("log_aggregations")

	doc := map[string]interface{}{
		"true":      body,
		"generated": *res,
	}

	_, err := coll.InsertOne(db.Ctx, doc)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func transformToTable(results []model.AggregateQuery, CountAggr bool, SumAggr bool, action string, origin string, brandId string, categoryId string) model.AggregateResult {
	cols := []string{"1m_bucket", "action"}
	// dynamically add columns in request
	if origin != "" {
		cols = append(cols, "origin")
	}
	if brandId != "" {
		cols = append(cols, "brand_id")
	}
	if categoryId != "" {
		cols = append(cols, "category_id")
	}
	// add aggregate columns
	if CountAggr {
		cols = append(cols, "count")
	}
	if SumAggr {
		cols = append(cols, "sum_price")
	}

	res := model.AggregateResult{Columns: cols}

	// iterate buckets
	for _, bucket := range results {
		if bucket.Id == "Other" {
			continue
		}

		// TODO make more efficent when init at the same time as cols
		// cut off timestring for bucket name
		bucketName := strings.TrimSuffix(bucket.Id, "Z")
		row := []string{bucketName, action}
		if origin != "" {
			row = append(row, origin)
		}
		if brandId != "" {
			row = append(row, brandId)
		}
		if categoryId != "" {
			row = append(row, categoryId)
		}

		if CountAggr {
			row = append(row, strconv.Itoa(bucket.Count))
		}
		if SumAggr {
			sumPrices := calcSumPrices(bucket.Tags)
			row = append(row, strconv.Itoa(sumPrices))
		}

		// append row to container
		res.Rows = append(res.Rows, row)
	}

	return res
}

func calcSumPrices(tags []model.AggregateTag) int {
	sum := 0
	for _, tag := range tags {
		sum += tag.ProductPrice
	}
	return sum
}
