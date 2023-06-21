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
		fmt.Println("DEBUG MODE")
	}
	return debug
}
