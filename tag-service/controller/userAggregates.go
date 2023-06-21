package controller

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	db "tag-service/database"
	model "tag-service/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func createTimeBoundaries(timestampFrom *time.Time, timestampEnd *time.Time) *[]time.Time {
	timeDiffMin := int(math.Ceil((*timestampEnd).Sub(*timestampFrom).Minutes()))
	buckets := []time.Time{}

	buckets = append(buckets, *timestampFrom)

	var t = *timestampFrom
	for i := 0; i < timeDiffMin-1; i++ {
		// one minute to each new bucket
		t = t.Add(time.Minute)
		buckets = append(buckets, t)
	}
	buckets = append(buckets, *timestampEnd)
	return &buckets
}

func queryDatabase(timestampFrom *time.Time, timestampEnd *time.Time, query *model.AggregateRequest) (*[]model.UserTagEvent, error) {
	q := *query

	filter := bson.A{
		bson.D{{"action", bson.D{{"$eq", q.Action}}}},
		bson.D{{"time", bson.D{{"$gte", primitive.NewDateTimeFromTime(*timestampFrom)}}}},
		bson.D{{"time", bson.D{{"$lte", primitive.NewDateTimeFromTime(*timestampEnd)}}}},
	}

	// add optional filter params
	if len(q.Origin) > 0 {
		d := bson.D{{"origin", bson.D{{"$eq", q.Origin}}}}
		filter = append(filter, d)
	}
	if len(q.CategoryId) > 0 {
		d := bson.D{{"productinfo.categoryid", bson.D{{"$eq", q.CategoryId}}}}
		filter = append(filter, d)
	}
	if len(q.BrandId) > 0 {
		d := bson.D{{"productinfo.brandid", bson.D{{"$eq", q.BrandId}}}}
		filter = append(filter, d)
	}

	filterArr := bson.D{
		{
			"$and",
			filter,
		},
	}

	// read user tags for cookie from database
	coll := db.DB.Database("mimuw").Collection("user_tags")

	// sort descending by timestamp
	opts := options.Find().SetSort(bson.D{{"time", 1}})

	cursor, err := coll.Find(db.Ctx, filterArr, opts)
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

func validateQuery(c *fiber.Ctx, debug bool) (*model.AggregateRequest, *time.Time, *time.Time, error) {
	// parse queries
	query := new(model.AggregateRequest)
	if err := c.QueryParser(query); err != nil {
		return nil, nil, nil, err
	}

	// print query params for debugging
	if debug {
		fmt.Println("time_range: ", query.TimeRange)
		fmt.Println("Action: ", query.Action)
		fmt.Println("Origin: ", query.Origin)
		fmt.Println("BrandId: ", query.BrandId)
		fmt.Println("CategoryId: ", query.CategoryId)
		fmt.Println("Aggregates: ", query.Aggregates)
	}

	// check time range query
	var timeRangeStr = query.TimeRange
	if timeRangeStr == "" {
		return nil, nil, nil, errors.New("Time range required")
	}
	// check and parse time range
	var timeRangeSplit = strings.Split(timeRangeStr, "_")

	// cast to time object
	var timeFormat = "2006-01-02T15:04:05" // exactly this format
	var timestampFrom, err0 = time.Parse(timeFormat, timeRangeSplit[0])
	var timestampEnd, err1 = time.Parse(timeFormat, timeRangeSplit[1])
	fmt.Println(timestampFrom, timestampEnd)
	if err0 != nil || err1 != nil {
		return nil, nil, nil, err0
	}

	// check action query
	var action = query.Action
	var validActions = []string{"BUY", "VIEW"}
	if action == "" || !contains(validActions, action) {
		return nil, nil, nil, errors.New("Valid action required")
	}

	// check aggregates query
	var validAggregates = []string{"COUNT", "SUM_PRICE"}
	if len(query.Aggregates) == 0 {
		return nil, nil, nil, errors.New("At least one aggregation required")
	}
	for _, aggr := range query.Aggregates {
		if !contains(validAggregates, aggr) {
			return nil, nil, nil, errors.New("Invalid aggregates")
		}
	}

	return query, &timestampFrom, &timestampEnd, nil
}

func initAggrMaps(buckets *[]time.Time) (map[string]int, map[string]int) {
	cnt_map := make(map[string]int)
	sum_map := make(map[string]int)
	for _, b := range *buckets {
		bStr := b.String()
		cnt_map[bStr] = 0
		sum_map[bStr] = 0
	}
	return cnt_map, sum_map
}

func createBucketTable(results *[]model.UserTagEvent, query *model.AggregateRequest, buckets *[]time.Time) *model.AggregateResult {
	q := *query
	useCnt := contains(q.Aggregates, "COUNT")
	useSum := contains(q.Aggregates, "SUM_PRICE")

	cols := []string{"1m_bucket", "action"}
	cnt_map, sum_map := initAggrMaps(buckets)

	// calculate count and sum_price for every minute bucket
	for _, res := range *results {
		// check to which time bucket the user tag belongs
		for _, bucket := range *buckets {
			resTime := res.Time.Truncate(time.Minute)
			if resTime == bucket.Truncate(time.Minute) {
				// add it to aggregates for this bucket
				cnt_map[bucket.String()] = cnt_map[bucket.String()] + 1
				sum_map[bucket.String()] = sum_map[bucket.String()] + res.ProductInfo.Price
			}
		}
	}

	// add optional params
	useOrigin := false
	if len(q.Origin) > 0 {
		cols = append(cols, "origin")
		useOrigin = true
	}
	useBrandId := false
	if len(q.BrandId) > 0 {
		cols = append(cols, "brand_id")
		useBrandId = true
	}
	useCategoryId := false
	if len(q.CategoryId) > 0 {
		cols = append(cols, "category_id")
		useCategoryId = true
	}
	// add aggregates
	if useCnt {
		cols = append(cols, "count")
	}
	if useSum {
		cols = append(cols, "sum_price")
	}

	// init rows with default values
	rows := [][]string{}

	for _, bucket := range *buckets {
		cnt := cnt_map[bucket.String()]
		price := sum_map[bucket.String()]

		// only include not empty buckets
		if cnt == 0 || price == 0 {
			continue
		}

		// Truncate the time to minutes
		trimmedTime := bucket.Truncate(time.Minute)
		// Format the time to the desired layout
		formattedTime := trimmedTime.Format("2006-01-02T15:04:05")

		row := []string{formattedTime, q.Action}
		if useOrigin {
			row = append(row, q.Origin)
		}
		if useBrandId {
			row = append(row, q.BrandId)
		}
		if useCategoryId {
			row = append(row, q.CategoryId)
		}

		// append aggregate results
		if useCnt {
			row = append(row, strconv.Itoa(cnt))
		}
		if useSum {
			row = append(row, strconv.Itoa(price))
		}

		// add to rows
		rows = append(rows, row)
	}

	// create aggregate result struct
	res := model.AggregateResult{Columns: cols, Rows: rows}
	return &res
}

func GetAggregate(c *fiber.Ctx, debug bool) error {
	// validate request query
	query, timestampFrom, timestampEnd, err := validateQuery(c, debug)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	results, err := queryDatabase(timestampFrom, timestampEnd, query)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// define 1min boundary timestamps for bucket aggregation
	buckets := createTimeBoundaries(timestampFrom, timestampEnd)

	// create 1 minute buckets out of results
	table := createBucketTable(results, query, buckets)

	// log response and expected result
	if debug {
		logAggrResponses(c, table, query)
	}

	return c.JSON(*table)
}

func logAggrResponses(c *fiber.Ctx, res *model.AggregateResult, query *model.AggregateRequest) {
	// debug response from api
	body := new(model.AggregateResult)
	c.BodyParser(&body)

	// ony log when generated response does not equal expected response
	areEqual := reflect.DeepEqual(body, *res)
	if areEqual {
		return
	}

	// actual generated response
	coll := db.DB.Database("mimuw").Collection("log_aggregations")

	doc := map[string]interface{}{
		"true":      body,
		"generated": *res,
		"query":     *query,
	}

	_, err := coll.InsertOne(db.Ctx, doc)
	if err != nil {
		fmt.Println(err.Error())
	}
}
