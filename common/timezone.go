package common

import (
	"log"
	"time"
)

var DhakaTZ *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		log.Fatalf("failed to load Dhaka timezone: %v", err)
	}
	DhakaTZ = loc
}
