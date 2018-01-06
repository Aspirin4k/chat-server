package utils

import (
	"io"
	"crypto/sha1"
)

func GetHash(input string) int {
	sha := sha1.New()
	io.WriteString(sha, input)
	return int(sha.Sum(nil)[0])
}