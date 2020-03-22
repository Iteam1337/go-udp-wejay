package utils

import (
	"fmt"
	"math"
	"time"
)

// RoundToNearestSecond …
func RoundToNearestSecond(now time.Time, interval time.Duration) time.Time {
	sec := now.Second() * int(time.Second)
	ns := now.Nanosecond()
	cur := sec + ns

	closest := math.Floor(float64(cur)/float64(interval)) * float64(interval)
	rounded := (int(closest) - cur) + int(interval)

	diff, _ := time.ParseDuration(fmt.Sprintf("%dns", rounded))

	return now.Add(diff)
}
