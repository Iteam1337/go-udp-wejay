package utils

import "math/rand"

// generate a number between >=0 && > max
func Random(max int) int {
	return int(rand.Float64() * float64(max))
}
