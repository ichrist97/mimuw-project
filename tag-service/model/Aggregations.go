package model

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
