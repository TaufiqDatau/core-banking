package util

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomInt(min, max int64) int64 {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return min + rng.Int63n(max-min+1)
}

func RandomBalance(min, max float64) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	randBalance := min + rng.Float64()*(max-min)

	return fmt.Sprintf("%.2f", randBalance)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString generates a random string of a given length
func RandomString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}
	return string(result)
}
