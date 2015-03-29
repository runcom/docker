package stringutils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	mathrand "math/rand"
	"strings"
	"time"
)

// Generate 32 chars random string
func GenerateRandomString() string {
	id := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		panic(err) // This shouldn't happen
	}
	return hex.EncodeToString(id)
}

// Generate alpha only random stirng with length n
func GenerateRandomAlphaOnlyString(n int) string {
	// make a really long string
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	r := mathrand.New(mathrand.NewSource(time.Now().UTC().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

// Generate Ascii random stirng with length n
func GenerateRandomAsciiString(n int) string {
	chars := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:` "
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = chars[mathrand.Intn(len(chars))]
	}
	return string(res)
}

// Truncate a string to maxlen
func Truncate(s string, maxlen int) string {
	if len(s) <= maxlen {
		return s
	}
	return s[:maxlen]
}

// Test wheather a string is contained in a slice of strings or not.
// Comparison is case insensitive
func InSlice(slice []string, s string) bool {
	for _, ss := range slice {
		if strings.ToLower(s) == strings.ToLower(ss) {
			return true
		}
	}
	return false
}
