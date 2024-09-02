package until

import (
	"math/rand"
	"strings"
	"time"
)

// Letters and digits to use for random string generation
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// init seeds the random number generator with the current time
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandString generates a random string of the given length
func RandString(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteRune(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}
