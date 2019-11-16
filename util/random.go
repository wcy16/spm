package util

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// generate random string with a random generator
func RandomStringWithGenerator(length int, r *rand.Rand) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

// generate random string
func RandomString(length int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return RandomStringWithGenerator(length, r)
}
