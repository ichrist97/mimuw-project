package model

import "time"

type AggregateRequest struct {
	TimeRange  string   `query:"time_range"`
	Action     string   `query:"action"`
	Origin     string   `query:"origin"`
	BrandId    string   `query:"brand_id"`
	CategoryId string   `query:"category_id"`
	Aggregates []string `query:"aggregates"`
}

type AggregateResult struct {
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
}

type AggregateTag struct {
	Action       string    `json:"action"`
	ProductPrice int       `json:"product_price" bson:"product_price"`
	Time         time.Time `json:"time"`
	Origin       string    `json:"origin,omitempty" bson:",omitempty"`
	BrandId      string    `json:"brand_id,omitempty" bson:",omitempty"`
	CategoryId   string    `json:"category_id,omitempty" bson:",omitempty"`
}
