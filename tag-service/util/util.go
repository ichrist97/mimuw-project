package util

import (
	"fmt"
	"os"
)

func DebugMode() bool {
	// activate debug mode
	debug_mode := os.Getenv("DEBUG")
	debug := false // default
	if debug_mode == "1" {
		debug = true
	}
	return debug
}

func AggregationMode() bool {
	// activate aggregations
	aggr_mode := os.Getenv("AGGR")
	aggr := true // default
	if aggr_mode == "0" {
		aggr = false
		fmt.Println("AGGREGATIONS DISABLED")
	}
	return aggr
}

func Contains(s []string, str string) bool {
	/**
	 * check if a string is contained in a list of strings
	 */
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
