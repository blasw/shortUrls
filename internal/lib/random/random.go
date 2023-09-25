package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var chs = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result []byte

	for i := 0; i < length; i++ {
		result = append(result, chs[rand.Intn(len(chs))])
	}

	return string(result)
}
