package util

import (
	"math/rand"
)

func GeneratorRandomString(n int) string {
	var charsets = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[rand.Intn(len(charsets))]
	}
	return string(letters)
}

func GeneratorRandomNumber(n int) string {
	var charsets = []rune("1234567890")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[rand.Intn(len(charsets))]
	}
	return string(letters)
}
