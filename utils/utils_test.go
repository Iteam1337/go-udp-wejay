package utils

import (
	"testing"
	"time"
)

func Test_RoundToNearestSecond(t *testing.T) {
	now := time.Now()
	interval := 5 * time.Second

	res := RoundToNearestSecond(now, interval)

	if res.Second()%5 != 0 {
		t.Error("expected second to be round down to closest duration (floor)", res.Second())
	}

	if res.Nanosecond() != 0 {
		t.Error("expected nanoseconds to be set to 0", res.Nanosecond())
	}
}
