package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
func RandomString(n int, allowedChars string) string {
	var letters = []rune(allowedChars)
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
